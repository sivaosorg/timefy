// Package main demonstrates an enterprise audit-trail system using timefy.
//
// Real-world scenario: a financial services platform must:
//
//  1. Record every data-access event with nanosecond-precision UTC timestamps.
//  2. Enforce a GDPR Article 33 compliance rule: any personal-data breach must
//     be reported to the supervisory authority within 72 hours of detection.
//  3. Detect SLA breaches (e.g., a support ticket must be resolved within
//     4 business hours of creation).
//  4. Generate a daily digest of events grouped by hour for the CISO dashboard.
//  5. Answer "how long ago was this event?" in human-readable form for the UI.
//
// Run with:
//
//	cd examples/audit && go run main.go
package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/sivaosorg/timefy"
)

// Severity classifies an audit event.
type Severity string

const (
	SeverityInfo     Severity = "INFO"
	SeverityWarning  Severity = "WARNING"
	SeverityCritical Severity = "CRITICAL"
	SeverityBreach   Severity = "BREACH"
)

// AuditEvent is a single immutable audit record.
type AuditEvent struct {
	ID         string
	OccurredAt time.Time // stored in UTC, nanosecond precision
	ActorID    string
	Action     string
	Resource   string
	Severity   Severity
	Resolved   bool
	ResolvedAt time.Time // zero value if unresolved
}

// HourBucket groups audit events that occurred within the same UTC hour.
type HourBucket struct {
	Hour   time.Time // beginning of the hour (UTC)
	Events []AuditEvent
}

// ─────────────────────────────────────────────────────────────────────────────
// gdprBreachDeadline returns the timestamp by which a detected breach must be
// reported.  GDPR Article 33 mandates notification within 72 hours.
//
// Why timefy.AddHour(detectedAt, 72)?
//   - Using time.Add(72 * time.Hour) is equivalent here, but AddHour reads as
//     intent-revealing code in an audit context where "hours" is the business
//     unit mandated by law.
//
// ─────────────────────────────────────────────────────────────────────────────
func gdprBreachDeadline(detectedAt time.Time) time.Time {
	return timefy.AddHour(detectedAt, 72)
}

// isGDPRCompliant reports whether a breach detected at `detectedAt` has been
// reported (i.e., is resolved) within the 72-hour window from `reportedAt`.
func isGDPRCompliant(detectedAt, reportedAt time.Time) bool {
	deadline := gdprBreachDeadline(detectedAt)
	// IsWithinTolerance would be too loose (±1 minute).  We use a direct
	// comparison: reported time must be before or exactly at the deadline.
	return !reportedAt.After(deadline)
}

// ─────────────────────────────────────────────────────────────────────────────
// slaDeadline returns the SLA resolution deadline for a ticket created at
// `createdAt`.  The SLA clock runs for `businessHours` business hours (M-F,
// 09:00-18:00 UTC for simplicity; adapt WithLocation for regional SLAs).
//
// Algorithm:
//  1. Start at createdAt.
//  2. Add one hour at a time, counting only hours that fall within business
//     hours on weekdays, until the required number of hours is consumed.
//
// Why not simply use AddHour(createdAt, businessHours)?
//   - SLAs in enterprise support contracts only count working hours.
//     A ticket created at 17:00 Friday with a 4-hour SLA should expire at
//     13:00 Monday, not 21:00 Friday (which would already be closed).
//
// ─────────────────────────────────────────────────────────────────────────────
func slaDeadline(createdAt time.Time, businessHours int) time.Time {
	cursor := createdAt
	remaining := businessHours

	for remaining > 0 {
		cursor = timefy.AddHour(cursor, 1)
		tx := timefy.New(cursor)
		h := cursor.Hour()
		if tx.IsWeekday() && h >= 9 && h < 18 {
			remaining--
		}
	}
	return cursor
}

// isSLABreached reports whether an event is unresolved and its SLA deadline
// has passed.
func isSLABreached(event AuditEvent, slaHours int) bool {
	if event.Resolved {
		return false
	}
	deadline := slaDeadline(event.OccurredAt, slaHours)
	return timefy.New(deadline).IsPast()
}

// ─────────────────────────────────────────────────────────────────────────────
// bucketByHour groups a slice of AuditEvents into HourBuckets sorted by the
// UTC hour in which they occurred.
//
// BeginningOfMinute (Truncate to minute) then BeginningOfHour gives us the
// canonical hour key.  Using timefy.New(e.OccurredAt).BeginningOfHour() is
// simpler and more readable than the equivalent time.Truncate(time.Hour).
// ─────────────────────────────────────────────────────────────────────────────
func bucketByHour(events []AuditEvent) []HourBucket {
	bucketMap := make(map[time.Time]*HourBucket)

	for _, e := range events {
		hourKey := timefy.New(e.OccurredAt).BeginningOfHour()
		if _, exists := bucketMap[hourKey]; !exists {
			bucketMap[hourKey] = &HourBucket{Hour: hourKey}
		}
		bucketMap[hourKey].Events = append(bucketMap[hourKey].Events, e)
	}

	// Sort buckets chronologically.
	buckets := make([]HourBucket, 0, len(bucketMap))
	for _, b := range bucketMap {
		buckets = append(buckets, *b)
	}
	sort.Slice(buckets, func(i, j int) bool {
		return buckets[i].Hour.Before(buckets[j].Hour)
	})
	return buckets
}

// ─────────────────────────────────────────────────────────────────────────────
// eventsInWindow returns all events whose OccurredAt falls within [start, end].
//
// IsBetween(start, end) performs an inclusive range check on the Timex, so
// events exactly at midnight or exactly at end-of-day are correctly included.
// ─────────────────────────────────────────────────────────────────────────────
func eventsInWindow(events []AuditEvent, start, end time.Time) []AuditEvent {
	var result []AuditEvent
	for _, e := range events {
		if timefy.New(e.OccurredAt).IsBetween(start, end) {
			result = append(result, e)
		}
	}
	return result
}

// ─────────────────────────────────────────────────────────────────────────────
// formatComponents uses FormatTimex to extract individual time components for
// a structured audit log line.
//
// FormatTimex returns [nanosecond, second, minute, hour, day, month, year]
// which maps perfectly to a structured log field set.
// ─────────────────────────────────────────────────────────────────────────────
func structuredTimestamp(t time.Time) string {
	c := timefy.FormatTimex(t)
	// c = [ns, sec, min, hour, day, month, year]
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%09d UTC",
		c[6], c[5], c[4], c[3], c[2], c[1], c[0])
}

func main() {
	// ── Synthetic audit log (timestamps carefully chosen to exercise each rule)

	// Reference "now" for this demo run.
	now := time.Date(2025, time.March, 24, 14, 30, 0, 0, time.UTC)

	events := []AuditEvent{
		// Normal access events today
		{
			ID: "evt-001", OccurredAt: time.Date(2025, 3, 24, 8, 15, 0, 0, time.UTC),
			ActorID: "alice", Action: "READ", Resource: "/api/users/42",
			Severity: SeverityInfo,
		},
		{
			ID: "evt-002", OccurredAt: time.Date(2025, 3, 24, 9, 5, 30, 123456789, time.UTC),
			ActorID: "bob", Action: "WRITE", Resource: "/api/payments/99",
			Severity: SeverityInfo,
		},
		// Warning: unusual bulk export
		{
			ID: "evt-003", OccurredAt: time.Date(2025, 3, 24, 10, 45, 0, 0, time.UTC),
			ActorID: "carol", Action: "EXPORT", Resource: "/api/users (bulk)",
			Severity: SeverityWarning,
		},
		// Critical: privilege escalation attempted
		{
			ID: "evt-004", OccurredAt: time.Date(2025, 3, 24, 11, 0, 0, 0, time.UTC),
			ActorID: "dave", Action: "ESCALATE", Resource: "/admin/roles",
			Severity: SeverityCritical, Resolved: false,
		},
		// BREACH: detected at 08:00, reported at 10:00 on the same day
		{
			ID: "evt-005", OccurredAt: time.Date(2025, 3, 22, 8, 0, 0, 0, time.UTC),
			ActorID: "system", Action: "DATA_LEAK", Resource: "/db/pii_records",
			Severity: SeverityBreach, Resolved: true,
			ResolvedAt: time.Date(2025, 3, 22, 10, 0, 0, 0, time.UTC),
		},
		// BREACH: detected Saturday, reported 80 hours later — NON-COMPLIANT
		{
			ID: "evt-006", OccurredAt: time.Date(2025, 3, 21, 6, 0, 0, 0, time.UTC),
			ActorID: "system", Action: "UNAUTHORIZED_ACCESS", Resource: "/api/credit_cards",
			Severity: SeverityBreach, Resolved: true,
			ResolvedAt: time.Date(2025, 3, 24, 14, 0, 0, 0, time.UTC), // 80 hours later
		},
		// Unresolved critical ticket — created Friday 17:30 UTC, SLA = 4 biz hours
		{
			ID: "evt-007", OccurredAt: time.Date(2025, 3, 21, 17, 30, 0, 0, time.UTC),
			ActorID: "frank", Action: "ACCESS_DENIED", Resource: "/api/reports",
			Severity: SeverityCritical, Resolved: false,
		},
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("      TIMEFY — ENTERPRISE AUDIT TRAIL DEMO               ")
	fmt.Println("═══════════════════════════════════════════════════════════")

	// ── 1. Full audit log with human-readable age ─────────────────────────────
	fmt.Println("\n  Full audit log:")
	fmt.Println("  ─────────────────────────────────────────────────────────────────────────")
	// We override "now" for TimeAgo by using DurationInDays / hours manually
	// when we need a frozen reference.  Standard TimeAgo() uses time.Now()
	// internally, which is correct for a live system.
	for _, e := range events {
		age := timefy.New(e.OccurredAt).TimeAgo()
		ts := structuredTimestamp(e.OccurredAt)
		fmt.Printf("  [%s]  %-12s  %-15s  %-10s  %s  (%s)\n",
			e.Severity, e.Action, e.Resource, e.ActorID, ts, age)
	}

	// ── 2. GDPR Article 33 compliance check ───────────────────────────────────
	fmt.Println("\n  GDPR Article 33 compliance (72-hour notification rule):")
	fmt.Println("  ─────────────────────────────────────────────────────────")
	for _, e := range events {
		if e.Severity != SeverityBreach {
			continue
		}
		deadline := gdprBreachDeadline(e.OccurredAt)
		if e.Resolved {
			compliant := isGDPRCompliant(e.OccurredAt, e.ResolvedAt)
			hoursUsed := int(e.ResolvedAt.Sub(e.OccurredAt).Hours())
			icon := "✓ COMPLIANT"
			if !compliant {
				icon = "✗ NON-COMPLIANT"
			}
			fmt.Printf("  [%s]  event %s  detected %s  deadline %s  reported in %dh  → %s\n",
				e.Severity, e.ID,
				e.OccurredAt.Format("01-02 15:04"),
				deadline.Format("01-02 15:04"),
				hoursUsed, icon)
		}
	}

	// ── 3. SLA breach detection ────────────────────────────────────────────────
	fmt.Println("\n  SLA breach detection (4 business hours):")
	fmt.Println("  ─────────────────────────────────────────────────────────")
	for _, e := range events {
		if e.Severity != SeverityCritical {
			continue
		}
		deadline := slaDeadline(e.OccurredAt, 4)
		breached := isSLABreached(e, 4)
		hoursElapsed := int(now.Sub(e.OccurredAt).Hours())
		icon := "✓ within SLA"
		if breached {
			icon = "✗ SLA BREACHED"
		}
		fmt.Printf("  event %-8s  created %s  SLA deadline %s  elapsed %dh  → %s\n",
			e.ID,
			e.OccurredAt.Format("Mon 01-02 15:04"),
			deadline.Format("Mon 01-02 15:04"),
			hoursElapsed, icon)
	}

	// ── 4. Events in today's window ───────────────────────────────────────────
	todayStart := timefy.New(now).BeginningOfDay()
	todayEnd := timefy.New(now).EndOfDay()
	todayEvents := eventsInWindow(events, todayStart, todayEnd)
	fmt.Printf("\n  Events on %s: %d total\n", now.Format("2006-01-02"), len(todayEvents))

	// ── 5. Hourly digest for today ────────────────────────────────────────────
	fmt.Println("\n  Hourly digest (today only):")
	fmt.Println("  ─────────────────────────────────────────────────────────")
	buckets := bucketByHour(todayEvents)
	for _, b := range buckets {
		fmt.Printf("  %s UTC — %d event(s)\n",
			b.Hour.Format("15:04"), len(b.Events))
		for _, e := range b.Events {
			fmt.Printf("    • [%s] %s %s by %s\n",
				e.Severity, e.Action, e.Resource, e.ActorID)
		}
	}

	// ── 6. Retention window: events older than 90 days can be archived ─────────
	fmt.Println("\n  Retention check (90-day archive threshold):")
	fmt.Println("  ─────────────────────────────────────────────────────────")
	archiveThreshold := timefy.BeginOfDay(timefy.AddDay(now, -90))
	for _, e := range events {
		if timefy.New(e.OccurredAt).IsBefore(archiveThreshold) {
			fmt.Printf("  ARCHIVE: %s  occurred %s\n",
				e.ID, e.OccurredAt.Format("2006-01-02"))
		}
	}
	fmt.Println("  (no events in this demo are older than 90 days)")

	// ── 7. Calendar context for the CISO weekly report ────────────────────────
	fmt.Println("\n  CISO weekly report context:")
	fmt.Println("  ─────────────────────────────────────────────────────────")
	tx := timefy.New(now)
	fmt.Printf("  Current week     : %s → %s\n",
		tx.BeginningOfWeek().Format("2006-01-02 Mon"),
		tx.EndOfWeek().Format("2006-01-02 Mon"))
	fmt.Printf("  Current quarter  : Q%d  %s → %s\n",
		tx.Quarter(),
		tx.BeginningOfQuarter().Format("2006-01-02"),
		tx.EndOfQuarter().Format("2006-01-02"))
	fmt.Printf("  Year-to-date     : %s → %s\n",
		tx.BeginningOfYear().Format("2006-01-02"),
		now.Format("2006-01-02"))
}
