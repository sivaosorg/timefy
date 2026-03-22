// Package main demonstrates enterprise-grade SaaS subscription billing period
// calculations using the timefy library.
//
// Real-world scenario: a multi-tenant SaaS platform must determine billing
// windows, pro-rate charges for mid-cycle sign-ups, detect trial expirations,
// and generate invoice line-items that align precisely with calendar boundaries.
//
// Run with:
//
//	cd _example/billing && go run main.go
package main

import (
	"fmt"
	"math"
	"time"

	"github.com/sivaosorg/timefy"
)

// Subscription holds all metadata for a single tenant's billing record.
type Subscription struct {
	TenantID    string
	Plan        string        // "monthly" | "quarterly" | "annual"
	StartedAt   time.Time     // exact moment the subscription became active
	PricePerDay float64       // daily rate in USD
	TrialDays   int           // number of free trial days (0 = no trial)
	Timezone    timefy.ZoneRFC // tenant's billing timezone
}

// InvoicePeriod describes a single billable window.
type InvoicePeriod struct {
	Start       time.Time
	End         time.Time
	BillableDays int
	Amount      float64
	Label       string
}

// ─────────────────────────────────────────────────────────────────────────────
// billingPeriodBoundaries returns the [start, end] of the billing period that
// contains `referenceTime` based on the subscription's plan and start date.
//
// Why timefy?
//   - BeginningOfMonth / EndOfMonth give exact calendar boundaries at
//     nanosecond precision, correctly handling month-length differences
//     (28/29/30/31 days) with a single call.
//   - BeginningOfQuarter / EndOfQuarter eliminate the common off-by-one month
//     bug when computing quarterly reports.
//   - BeginningOfYear / EndOfYear handle the year-wrap edge case.
// ─────────────────────────────────────────────────────────────────────────────
func billingPeriodBoundaries(s Subscription, referenceTime time.Time) (start, end time.Time) {
	// Convert reference time to the tenant's billing timezone so that
	// month and year boundaries are evaluated in their local context.
	localRef := timefy.AdjustTimezone(referenceTime, string(s.Timezone))
	tx := timefy.New(localRef)

	switch s.Plan {
	case "monthly":
		// The billing window is the calendar month containing localRef.
		start = tx.BeginningOfMonth()
		end = tx.EndOfMonth()
	case "quarterly":
		// Q1 = Jan-Mar, Q2 = Apr-Jun, Q3 = Jul-Sep, Q4 = Oct-Dec.
		start = tx.BeginningOfQuarter()
		end = tx.EndOfQuarter()
	case "annual":
		// Calendar year January 1 → December 31.
		start = tx.BeginningOfYear()
		end = tx.EndOfYear()
	default:
		// Fallback: 30-day window anchored to the subscription start day.
		start = timefy.BeginOfDay(s.StartedAt)
		end = timefy.FEndOfDay(timefy.AddDay(start, 30))
	}
	return
}

// ─────────────────────────────────────────────────────────────────────────────
// trialEndsAt returns the exact moment the free trial expires.
//
// Why timefy.AddDay?
//   - timefy.AddDay wraps time.Add(24 * n * time.Hour), which correctly
//     accounts for DST transitions — unlike AddDate(0, 0, n) when the tenant
//     timezone observes daylight saving.
// ─────────────────────────────────────────────────────────────────────────────
func trialEndsAt(s Subscription) time.Time {
	if s.TrialDays == 0 {
		return s.StartedAt
	}
	trialEnd := timefy.AddDay(s.StartedAt, s.TrialDays)
	// Expire at the last nanosecond of the trial's final day.
	return timefy.New(trialEnd).EndOfDay()
}

// ─────────────────────────────────────────────────────────────────────────────
// isInTrial reports whether `now` falls within the subscription's trial window.
// ─────────────────────────────────────────────────────────────────────────────
func isInTrial(s Subscription, now time.Time) bool {
	if s.TrialDays == 0 {
		return false
	}
	trialEnd := trialEndsAt(s)
	// IsBetween is inclusive on both ends, which is what billing systems need:
	// the trial is active on day 1 through the last nanosecond of day N.
	return timefy.New(now).IsBetween(s.StartedAt, trialEnd)
}

// ─────────────────────────────────────────────────────────────────────────────
// proRateForPeriod calculates a pro-rated charge when the subscription started
// mid-period.  The function determines how many days of the billing period were
// actually used and charges accordingly.
//
// Example: Annual plan, $1 / day.  Tenant signed up on March 15.
//   - Annual period: Jan 1 – Dec 31 (365 days).
//   - Days actually used: March 15 → Dec 31 = 291 days.
//   - Charge: 291 × $1 = $291 (instead of $365).
// ─────────────────────────────────────────────────────────────────────────────
func proRateForPeriod(s Subscription, periodStart, periodEnd time.Time) InvoicePeriod {
	// Effective start is the later of periodStart and the subscription start.
	// This handles mid-cycle sign-ups correctly.
	effectiveStart := periodStart
	if s.StartedAt.After(periodStart) {
		effectiveStart = timefy.BeginOfDay(s.StartedAt)
	}

	// Effective end is the earlier of periodEnd and today's end-of-day
	// (we don't bill for future days).
	effectiveEnd := periodEnd
	now := time.Now()
	if now.Before(periodEnd) {
		effectiveEnd = timefy.New(now).EndOfDay()
	}

	// DurationInDays uses (t.Sub(other).Hours() / 24) which gives calendar days.
	// We add 1 because both boundaries are inclusive.
	days := int(math.Abs(float64(timefy.New(effectiveEnd).DurationInDays(effectiveStart)))) + 1

	// Deduct trial days that fall within this period.
	if s.TrialDays > 0 {
		trialEnd := trialEndsAt(s)
		if trialEnd.After(effectiveStart) {
			trialInPeriod := int(math.Abs(float64(timefy.New(trialEnd).DurationInDays(effectiveStart)))) + 1
			days -= trialInPeriod
			if days < 0 {
				days = 0
			}
		}
	}

	amount := float64(days) * s.PricePerDay

	return InvoicePeriod{
		Start:        effectiveStart,
		End:          effectiveEnd,
		BillableDays: days,
		Amount:       amount,
		Label: fmt.Sprintf("%s → %s",
			effectiveStart.Format("2006-01-02"),
			effectiveEnd.Format("2006-01-02"),
		),
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// generateUpcomingInvoices produces the next N invoice periods starting from
// `from`.  This powers an invoice preview feature on a customer dashboard.
// ─────────────────────────────────────────────────────────────────────────────
func generateUpcomingInvoices(s Subscription, from time.Time, count int) []InvoicePeriod {
	invoices := make([]InvoicePeriod, 0, count)
	cursor := from

	for len(invoices) < count {
		start, end := billingPeriodBoundaries(s, cursor)
		// Move to the next period: advance cursor by one day past the period end.
		cursor = timefy.AddDay(end, 1)

		days := int(math.Abs(float64(timefy.New(end).DurationInDays(start)))) + 1
		invoices = append(invoices, InvoicePeriod{
			Start:        start,
			End:          end,
			BillableDays: days,
			Amount:       float64(days) * s.PricePerDay,
			Label:        fmt.Sprintf("%s → %s", start.Format("2006-01-02"), end.Format("2006-01-02")),
		})
	}
	return invoices
}

func main() {
	// ── Tenant A: monthly plan, started on the 15th → pro-rated first month ──
	tenantA := Subscription{
		TenantID:    "tenant-a",
		Plan:        "monthly",
		StartedAt:   time.Date(2025, time.March, 15, 9, 0, 0, 0, time.UTC),
		PricePerDay: 3.33, // ~$100 / month
		TrialDays:   0,
		Timezone:    timefy.DefaultTimezoneNewYork,
	}

	// ── Tenant B: annual plan, 14-day trial, started February 1 ──
	tenantB := Subscription{
		TenantID:    "tenant-b",
		Plan:        "annual",
		StartedAt:   time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC),
		PricePerDay: 1.00, // $365 / year
		TrialDays:   14,
		Timezone:    timefy.DefaultTimezoneVietnam,
	}

	// ── Tenant C: quarterly plan ──
	tenantC := Subscription{
		TenantID:    "tenant-c",
		Plan:        "quarterly",
		StartedAt:   time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
		PricePerDay: 3.33,
		Timezone:    timefy.DefaultTimezoneLondon,
	}

	referenceNow := time.Date(2025, time.March, 22, 12, 0, 0, 0, time.UTC)

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("          TIMEFY — ENTERPRISE BILLING DEMO                ")
	fmt.Println("═══════════════════════════════════════════════════════════")

	// ── Tenant A: current billing window ──
	fmt.Println("\n── Tenant A (monthly, no trial) ────────────────────────────")
	startA, endA := billingPeriodBoundaries(tenantA, referenceNow)
	fmt.Printf("  Current billing period : %s → %s\n",
		startA.Format("2006-01-02 15:04:05 MST"),
		endA.Format("2006-01-02 15:04:05 MST"),
	)
	invoiceA := proRateForPeriod(tenantA, startA, endA)
	fmt.Printf("  Pro-rated invoice      : %s  (%d days × $%.2f = $%.2f)\n",
		invoiceA.Label, invoiceA.BillableDays, tenantA.PricePerDay, invoiceA.Amount)

	// ── Tenant B: trial status + annual period ──
	fmt.Println("\n── Tenant B (annual, 14-day trial) ─────────────────────────")
	trialEnd := trialEndsAt(tenantB)
	fmt.Printf("  Trial ends at          : %s\n", trialEnd.Format("2006-01-02 15:04:05"))
	fmt.Printf("  In trial now?          : %v\n", isInTrial(tenantB, referenceNow))
	startB, endB := billingPeriodBoundaries(tenantB, referenceNow)
	invoiceB := proRateForPeriod(tenantB, startB, endB)
	fmt.Printf("  Annual period          : %s → %s\n",
		startB.Format("2006-01-02"), endB.Format("2006-01-02"))
	fmt.Printf("  Billable days (YTD)    : %d  →  $%.2f\n",
		invoiceB.BillableDays, invoiceB.Amount)

	// ── Tenant C: upcoming quarterly invoices ──
	fmt.Println("\n── Tenant C (quarterly) — next 4 invoices ──────────────────")
	upcoming := generateUpcomingInvoices(tenantC, referenceNow, 4)
	for i, inv := range upcoming {
		fmt.Printf("  Invoice %d: %s  (%d days, $%.2f)\n",
			i+1, inv.Label, inv.BillableDays, inv.Amount)
	}

	// ── Demonstrate leap-year edge-case ──────────────────────────────────────
	fmt.Println("\n── Leap-year February edge-case ─────────────────────────────")
	leap2024 := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	leapMonth := timefy.New(leap2024).EndOfMonth()
	fmt.Printf("  End of February 2024   : %s (leap: %v)\n",
		leapMonth.Format("2006-01-02"), timefy.IsLeapYear(2024))
	nonLeap2023 := time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC)
	nonLeapMonth := timefy.New(nonLeap2023).EndOfMonth()
	fmt.Printf("  End of February 2023   : %s (leap: %v)\n",
		nonLeapMonth.Format("2006-01-02"), timefy.IsLeapYear(2023))

	// ── Demonstrate timezone-aware monthly boundary ───────────────────────────
	fmt.Println("\n── Timezone-aware billing boundary demo ─────────────────────")
	// 2025-03-31 23:00 UTC is 2025-04-01 09:00 in Vietnam (+7).
	// For a Vietnamese tenant on a monthly plan, this moment belongs to April,
	// not March — timefy.AdjustTimezone ensures the boundary is evaluated
	// in the correct timezone.
	borderTime := time.Date(2025, time.March, 31, 23, 0, 0, 0, time.UTC)
	vnTime := timefy.AdjustTimezone(borderTime, string(timefy.DefaultTimezoneVietnam))
	vnTx := timefy.New(vnTime)
	fmt.Printf("  UTC time               : %s\n", borderTime.Format("2006-01-02 15:04 MST"))
	fmt.Printf("  Vietnam time           : %s\n", vnTime.Format("2006-01-02 15:04 MST"))
	fmt.Printf("  Vietnam billing month  : %s → %s\n",
		vnTx.BeginningOfMonth().Format("2006-01-02"),
		vnTx.EndOfMonth().Format("2006-01-02"))
}
