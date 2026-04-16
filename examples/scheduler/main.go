// Package main demonstrates a business-day-aware job scheduler using timefy.
//
// Real-world scenario: a financial data platform runs a suite of ETL jobs
// that must execute on specific schedules:
//
//  1. Daily market-data ingestion: runs at 18:30 every weekday (markets closed).
//  2. End-of-month portfolio reconciliation: runs at 02:00 on the first
//     business day of the month after month-end.
//  3. Quarterly regulatory report: runs at 08:00 on the first weekday of the
//     new quarter.
//  4. Rate-limiting guard: no job may be triggered more than once per
//     configurable time window (e.g. 5 minutes).
//  5. Job-run history with human-readable "last ran" display.
//
// Run with:
//
//	cd examples/scheduler && go run main.go
package main

import (
	"fmt"
	"time"

	"github.com/sivaosorg/timefy"
)

// Schedule describes a recurring job's timing rules.
type Schedule struct {
	Name            string
	RunHour         int           // UTC hour for the job to fire
	RunMinute       int           // UTC minute for the job to fire
	Frequency       string        // "daily" | "monthly" | "quarterly"
	WindowMinutes   int           // anti-duplication window in minutes
	Timezone        timefy.ZoneRFC // timezone for schedule evaluation
}

// JobRun records a single execution of a scheduled job.
type JobRun struct {
	JobName   string
	RanAt     time.Time
	DurationMs int64
	Status    string // "success" | "failure"
}

// ─────────────────────────────────────────────────────────────────────────────
// nextRunTime computes the next UTC timestamp at which a job should fire,
// based on the schedule and a reference "now" time.
//
// For daily jobs:
//   - If today's scheduled time is still in the future, return it.
//   - Otherwise, return the same time on the next weekday.
//
// For monthly jobs:
//   - Return the first weekday of the next month at the configured hour.
//   - "First weekday" means we step forward from the 1st until IsWeekday().
//
// For quarterly jobs:
//   - Return the first weekday of the next quarter at the configured hour.
//
// Why timefy?
//   - BeginningOfMonth / BeginningOfQuarter give exact first-of-period
//     boundaries regardless of month length or DST.
//   - AddDay + IsWeekday allow clean forward-stepping without manual
//     month-wrap arithmetic.
// ─────────────────────────────────────────────────────────────────────────────
func nextRunTime(s Schedule, now time.Time) time.Time {
	// Evaluate schedule in the job's configured timezone.
	localNow := timefy.AdjustTimezone(now, string(s.Timezone))

	switch s.Frequency {
	case "daily":
		return nextDailyRun(s, localNow)
	case "monthly":
		return nextMonthlyRun(s, localNow)
	case "quarterly":
		return nextQuarterlyRun(s, localNow)
	default:
		return time.Time{}
	}
}

func nextDailyRun(s Schedule, localNow time.Time) time.Time {
	// Candidate: today at the scheduled local time.
	todayRun := time.Date(
		localNow.Year(), localNow.Month(), localNow.Day(),
		s.RunHour, s.RunMinute, 0, 0, localNow.Location(),
	)

	cursor := todayRun
	// If today's run is in the past, advance to tomorrow and keep stepping
	// until we land on a weekday.
	if !cursor.After(localNow) {
		cursor = timefy.AddDay(cursor, 1)
	}
	for !timefy.New(cursor).IsWeekday() {
		cursor = timefy.AddDay(cursor, 1)
	}
	return cursor.UTC()
}

func nextMonthlyRun(s Schedule, localNow time.Time) time.Time {
	// Start looking from the beginning of next month.
	tx := timefy.New(localNow)
	nextMonthFirst := tx.BeginningOfMonth().AddDate(0, 1, 0)

	// Step forward from the 1st until we find a weekday.
	cursor := time.Date(
		nextMonthFirst.Year(), nextMonthFirst.Month(), nextMonthFirst.Day(),
		s.RunHour, s.RunMinute, 0, 0, nextMonthFirst.Location(),
	)
	for !timefy.New(cursor).IsWeekday() {
		cursor = timefy.AddDay(cursor, 1)
	}
	return cursor.UTC()
}

func nextQuarterlyRun(s Schedule, localNow time.Time) time.Time {
	// Start looking from the beginning of the next calendar quarter.
	tx := timefy.New(localNow)
	nextQFirst := tx.BeginningOfQuarter().AddDate(0, 3, 0)

	cursor := time.Date(
		nextQFirst.Year(), nextQFirst.Month(), nextQFirst.Day(),
		s.RunHour, s.RunMinute, 0, 0, nextQFirst.Location(),
	)
	for !timefy.New(cursor).IsWeekday() {
		cursor = timefy.AddDay(cursor, 1)
	}
	return cursor.UTC()
}

// ─────────────────────────────────────────────────────────────────────────────
// isDuplicateRun returns true if a job with the same name already ran within
// the anti-duplication window (in minutes) before `now`.
//
// The check uses IsBetween to determine whether any prior run falls inside
// the [now − window, now] interval.
//
// Why IsBetween over a raw time.Sub comparison?
//   - IsBetween is inclusive on both ends (equal to start or end is "within"),
//     which correctly prevents a run at exactly the window boundary.
// ─────────────────────────────────────────────────────────────────────────────
func isDuplicateRun(jobName string, history []JobRun, windowMinutes int, now time.Time) bool {
	windowStart := timefy.AddMinute(now, -windowMinutes)
	for _, run := range history {
		if run.JobName == jobName && timefy.New(run.RanAt).IsBetween(windowStart, now) {
			return true
		}
	}
	return false
}

// ─────────────────────────────────────────────────────────────────────────────
// businessDaysUntilRun computes how many weekdays remain between `now` and
// the job's next run time.  This is displayed on the ops dashboard as a
// human-readable countdown.
// ─────────────────────────────────────────────────────────────────────────────
func businessDaysUntilRun(now, nextRun time.Time) int {
	weekdays := timefy.GetWeekdaysInRange(timefy.BeginOfDay(now), timefy.BeginOfDay(nextRun))
	return len(weekdays)
}

// ─────────────────────────────────────────────────────────────────────────────
// renderJobCard prints a formatted summary for a single scheduled job.
// ─────────────────────────────────────────────────────────────────────────────
func renderJobCard(s Schedule, lastRun *JobRun, now time.Time) {
	next := nextRunTime(s, now)
	bizDays := businessDaysUntilRun(now, next)

	fmt.Printf("\n  ┌─ %-40s ─────────────────\n", s.Name)
	fmt.Printf("  │  Frequency         : %s\n", s.Frequency)
	fmt.Printf("  │  Next run (UTC)     : %s\n", next.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("  │  Business days away : %d\n", bizDays)
	fmt.Printf("  │  Countdown          : %s\n", timefy.New(next).TimeUntil())
	if lastRun != nil {
		fmt.Printf("  │  Last ran           : %s  [%s]  (%s)\n",
			lastRun.RanAt.Format("2006-01-02 15:04 MST"),
			lastRun.Status,
			timefy.New(lastRun.RanAt).TimeAgo(),
		)
	} else {
		fmt.Printf("  │  Last ran           : never\n")
	}
	fmt.Printf("  └──────────────────────────────────────────────────────────\n")
}

func main() {
	// Reference "now" for deterministic demo output.
	now := time.Date(2025, time.March, 24, 14, 0, 0, 0, time.UTC)

	// ── Job definitions ───────────────────────────────────────────────────────
	jobs := []Schedule{
		{
			Name:          "market-data-ingestion",
			RunHour:       18,
			RunMinute:     30,
			Frequency:     "daily",
			WindowMinutes: 5,
			Timezone:      timefy.DefaultTimezoneNewYork,
		},
		{
			Name:          "monthly-portfolio-reconciliation",
			RunHour:       2,
			RunMinute:     0,
			Frequency:     "monthly",
			WindowMinutes: 60,
			Timezone:      timefy.DefaultTimezoneLondon,
		},
		{
			Name:          "quarterly-regulatory-report",
			RunHour:       8,
			RunMinute:     0,
			Frequency:     "quarterly",
			WindowMinutes: 120,
			Timezone:      timefy.DefaultTimezoneNewYork,
		},
	}

	// ── Historical run log (simulate some prior executions) ───────────────────
	history := []JobRun{
		{
			JobName:    "market-data-ingestion",
			RanAt:      time.Date(2025, 3, 21, 23, 30, 0, 0, time.UTC),
			DurationMs: 4312,
			Status:     "success",
		},
		{
			JobName:    "monthly-portfolio-reconciliation",
			RanAt:      time.Date(2025, 3, 3, 2, 0, 0, 0, time.UTC),
			DurationMs: 18200,
			Status:     "success",
		},
		{
			JobName:    "quarterly-regulatory-report",
			RanAt:      time.Date(2025, 1, 2, 8, 0, 0, 0, time.UTC),
			DurationMs: 92100,
			Status:     "success",
		},
	}

	// ── Simulate a duplicate-trigger attempt ──────────────────────────────────
	recentHistory := []JobRun{
		{
			JobName: "market-data-ingestion",
			RanAt:   timefy.AddMinute(now, -2), // ran 2 minutes ago
			Status:  "success",
		},
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("     TIMEFY — BUSINESS-DAY-AWARE JOB SCHEDULER DEMO     ")
	fmt.Println("═══════════════════════════════════════════════════════════")

	// ── Render job cards ──────────────────────────────────────────────────────
	for _, job := range jobs {
		var last *JobRun
		for i := range history {
			if history[i].JobName == job.Name {
				lr := history[i]
				last = &lr
				break
			}
		}
		renderJobCard(job, last, now)
	}

	// ── Anti-duplication check ────────────────────────────────────────────────
	fmt.Println("\n  Anti-duplication window checks (5 min):")
	fmt.Println("  ─────────────────────────────────────────────────────────")
	for _, job := range jobs {
		dup := isDuplicateRun(job.Name, recentHistory, job.WindowMinutes, now)
		icon := "✓ OK to run"
		if dup {
			icon = "✗ DUPLICATE — suppressed"
		}
		fmt.Printf("  %-40s → %s\n", job.Name, icon)
	}

	// ── Weekly schedule overview (Mon–Fri) ────────────────────────────────────
	fmt.Println("\n  Daily job — next 5 occurrences:")
	fmt.Println("  ─────────────────────────────────────────────────────────")
	dailyJob := jobs[0]
	cursor := now
	for i := 0; i < 5; i++ {
		next := nextRunTime(dailyJob, cursor)
		fmt.Printf("  Run %d: %s UTC\n", i+1, next.Format("2006-01-02 (Mon) 15:04"))
		// Advance past this run so the next iteration finds the following day.
		cursor = timefy.AddMinute(next, 1)
	}

	// ── Month boundary: first weekday of each upcoming month ──────────────────
	fmt.Println("\n  Monthly job — next 3 first-weekday-of-month dates:")
	fmt.Println("  ─────────────────────────────────────────────────────────")
	monthlyJob := jobs[1]
	mCursor := now
	for i := 0; i < 3; i++ {
		next := nextRunTime(monthlyJob, mCursor)
		fmt.Printf("  %s UTC\n", next.Format("2006-01-02 (Mon) 15:04"))
		mCursor = timefy.AddDay(next, 1)
	}

	// ── Quarter boundary: first weekday of each upcoming quarter ──────────────
	fmt.Println("\n  Quarterly job — next 4 first-weekday-of-quarter dates:")
	fmt.Println("  ─────────────────────────────────────────────────────────")
	quarterlyJob := jobs[2]
	qCursor := now
	for i := 0; i < 4; i++ {
		next := nextRunTime(quarterlyJob, qCursor)
		fmt.Printf("  %s UTC  (Q%d)\n",
			next.Format("2006-01-02 (Mon) 15:04"),
			timefy.New(next).Quarter(),
		)
		qCursor = timefy.AddDay(next, 1)
	}

	// ── Time-window calculations ───────────────────────────────────────────────
	fmt.Println("\n  Utility calculations:")
	fmt.Println("  ─────────────────────────────────────────────────────────")
	fmt.Printf("  Current quarter   : Q%d\n", timefy.New(now).Quarter())
	fmt.Printf("  Beginning of week : %s\n", timefy.New(now).BeginningOfWeek().Format("2006-01-02 Mon"))
	fmt.Printf("  End of week       : %s\n", timefy.New(now).EndOfWeek().Format("2006-01-02 Mon"))
	fmt.Printf("  Weekdays this week:\n")
	weekStart := timefy.New(now).BeginningOfWeek()
	weekEnd := timefy.New(now).EndOfWeek()
	for _, d := range timefy.GetWeekdaysInRange(weekStart, weekEnd) {
		fmt.Printf("    • %s\n", d.Format("Mon 2006-01-02"))
	}
}
