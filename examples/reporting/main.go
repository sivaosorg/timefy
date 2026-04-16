// Package main demonstrates financial quarter and fiscal-year reporting using
// the timefy library.
//
// Real-world scenario: a publicly-traded technology company must produce:
//
//  1. SEC quarterly filing windows (10-Q: 45 days after quarter end).
//  2. Annual report window (10-K: 60 days after fiscal year end).
//  3. Year-to-date (YTD) revenue aggregation boundaries.
//  4. Quarter-over-quarter (QoQ) and year-over-year (YoY) comparison periods.
//  5. A fiscal-year calendar that starts on July 1 (non-calendar fiscal year).
//
// Run with:
//
//	cd _example/reporting && go run main.go
package main

import (
	"fmt"
	"time"

	"github.com/sivaosorg/timefy"
)

// QuarterPeriod holds the boundaries and metadata for a single reporting quarter.
type QuarterPeriod struct {
	Label    string    // e.g. "Q3 FY2025"
	Start    time.Time
	End      time.Time
	FilingDeadline time.Time // SEC 10-Q deadline
	Quarter  uint
	FiscalYear int
}

// FiscalYear holds boundaries for a full fiscal year (non-calendar).
type FiscalYearPeriod struct {
	Label          string
	Start          time.Time
	End            time.Time
	FilingDeadline time.Time // SEC 10-K deadline
	FiscalYear     int
}

// ─────────────────────────────────────────────────────────────────────────────
// calendarQuarterPeriod returns the QuarterPeriod for the calendar quarter that
// contains `t`.
//
// BeginningOfQuarter and EndOfQuarter are the two key functions here.
// They correctly compute quarter boundaries without any switch-case month
// arithmetic: the library handles "which 3-month block does this month
// belong to?" internally.
// ─────────────────────────────────────────────────────────────────────────────
func calendarQuarterPeriod(t time.Time) QuarterPeriod {
	tx := timefy.New(t)
	start := tx.BeginningOfQuarter()
	end := tx.EndOfQuarter()
	q := tx.Quarter()
	year := t.Year()

	// SEC 10-Q filing deadline: 45 days after quarter end for large accelerated
	// filers.  We add 45 calendar days to the end of quarter.
	filingDeadline := timefy.New(end).EndOfDay()
	filingDeadline = timefy.AddDay(filingDeadline, 45)

	return QuarterPeriod{
		Label:          fmt.Sprintf("Q%d CY%d", q, year),
		Start:          start,
		End:            end,
		FilingDeadline: filingDeadline,
		Quarter:        q,
		FiscalYear:     year,
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// fiscalQuarterPeriod returns the QuarterPeriod for a July-1-start fiscal year.
//
// Fiscal year convention used here (common in technology companies):
//   FY2025 = July 1 2024 → June 30 2025
//   FQ1 FY2025 = Jul 1 – Sep 30 2024
//   FQ2 FY2025 = Oct 1 – Dec 31 2024
//   FQ3 FY2025 = Jan 1 – Mar 31 2025
//   FQ4 FY2025 = Apr 1 – Jun 30 2025
//
// Because timefy operates on calendar months, we adjust the reference month
// by −6 to align with the fiscal-year start before calling BeginningOfQuarter.
// ─────────────────────────────────────────────────────────────────────────────
func fiscalQuarterPeriod(t time.Time) QuarterPeriod {
	// Shift reference time 6 months backward so that the standard calendar
	// quarter functions align with a July-start fiscal year.
	shifted := t.AddDate(0, -6, 0)
	tx := timefy.New(shifted)

	// Calendar quarter of the shifted date maps to the fiscal quarter.
	calStart := tx.BeginningOfQuarter()
	calEnd := tx.EndOfQuarter()

	// Shift boundaries back to real dates.
	fqStart := calStart.AddDate(0, 6, 0)
	fqEnd := calEnd.AddDate(0, 6, 0)

	fq := tx.Quarter()

	// Fiscal year label: the year in which June 30 falls.
	fy := fqStart.Year()
	if fqStart.Month() >= time.July {
		fy++ // FY label increments after July 1
	}

	filingDeadline := timefy.AddDay(timefy.New(fqEnd).EndOfDay(), 45)

	return QuarterPeriod{
		Label:          fmt.Sprintf("FQ%d FY%d", fq, fy),
		Start:          fqStart,
		End:            fqEnd,
		FilingDeadline: filingDeadline,
		Quarter:        fq,
		FiscalYear:     fy,
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// fiscalYearPeriod returns the full fiscal year period that contains `t`.
// FY2025 runs July 1 2024 → June 30 2025.
// ─────────────────────────────────────────────────────────────────────────────
func fiscalYearPeriod(t time.Time) FiscalYearPeriod {
	shifted := t.AddDate(0, -6, 0)
	tx := timefy.New(shifted)

	calYearStart := tx.BeginningOfYear()
	calYearEnd := tx.EndOfYear()

	fyStart := calYearStart.AddDate(0, 6, 0)
	fyEnd := calYearEnd.AddDate(0, 6, 0)

	fy := fyStart.Year()
	if fyStart.Month() >= time.July {
		fy++
	}

	// SEC 10-K: 60 days after fiscal year end for large accelerated filers.
	filingDeadline := timefy.AddDay(timefy.New(fyEnd).EndOfDay(), 60)

	return FiscalYearPeriod{
		Label:          fmt.Sprintf("FY%d", fy),
		Start:          fyStart,
		End:            fyEnd,
		FilingDeadline: filingDeadline,
		FiscalYear:     fy,
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// ytdPeriod returns the year-to-date window: from January 1 of the current
// calendar year to `referenceDate` (inclusive).
//
// BeginningOfYear gives January 1 00:00:00 and EndOfDay gives the last
// nanosecond of `referenceDate`, forming the exact YTD window used by finance
// teams when querying a data warehouse.
// ─────────────────────────────────────────────────────────────────────────────
func ytdPeriod(referenceDate time.Time) (start, end time.Time) {
	start = timefy.New(referenceDate).BeginningOfYear()
	end = timefy.New(referenceDate).EndOfDay()
	return
}

// ─────────────────────────────────────────────────────────────────────────────
// qoqComparison returns the current quarter and the same quarter one year ago,
// which is the standard quarter-over-quarter comparison in financial analysis.
// ─────────────────────────────────────────────────────────────────────────────
func qoqComparison(referenceDate time.Time) (current, priorYear QuarterPeriod) {
	current = calendarQuarterPeriod(referenceDate)
	priorYear = calendarQuarterPeriod(referenceDate.AddDate(-1, 0, 0))
	return
}

// ─────────────────────────────────────────────────────────────────────────────
// daysTillFiling returns a human-readable countdown to a filing deadline.
//
// DurationInDays computes the absolute difference in whole days between two
// Timex instances.  We use it here as a countdown timer.
// ─────────────────────────────────────────────────────────────────────────────
func daysTillFiling(deadline, referenceDate time.Time) int {
	tx := timefy.New(deadline)
	days := tx.DurationInDays(referenceDate)
	if days < 0 {
		return 0
	}
	return days
}

func main() {
	referenceDate := time.Date(2025, time.March, 22, 12, 0, 0, 0, time.UTC)

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("      TIMEFY — FINANCIAL REPORTING CALENDAR DEMO         ")
	fmt.Println("═══════════════════════════════════════════════════════════")

	// ── Calendar quarter ────────────────────────────────────────────────────
	cq := calendarQuarterPeriod(referenceDate)
	fmt.Printf("\n  Calendar quarter for %s:\n", referenceDate.Format("2006-01-02"))
	fmt.Printf("  %-12s  %s → %s\n", cq.Label,
		cq.Start.Format("2006-01-02"), cq.End.Format("2006-01-02"))
	fmt.Printf("  10-Q filing deadline : %s  (%d days away)\n",
		cq.FilingDeadline.Format("2006-01-02"),
		daysTillFiling(cq.FilingDeadline, referenceDate))

	// ── Fiscal quarter (July-start FY) ──────────────────────────────────────
	fq := fiscalQuarterPeriod(referenceDate)
	fmt.Printf("\n  Fiscal quarter (Jul-start) for %s:\n", referenceDate.Format("2006-01-02"))
	fmt.Printf("  %-12s  %s → %s\n", fq.Label,
		fq.Start.Format("2006-01-02"), fq.End.Format("2006-01-02"))
	fmt.Printf("  10-Q filing deadline : %s  (%d days away)\n",
		fq.FilingDeadline.Format("2006-01-02"),
		daysTillFiling(fq.FilingDeadline, referenceDate))

	// ── Fiscal year ─────────────────────────────────────────────────────────
	fy := fiscalYearPeriod(referenceDate)
	fmt.Printf("\n  %s  %s → %s\n", fy.Label,
		fy.Start.Format("2006-01-02"), fy.End.Format("2006-01-02"))
	fmt.Printf("  10-K filing deadline : %s  (%d days away)\n",
		fy.FilingDeadline.Format("2006-01-02"),
		daysTillFiling(fy.FilingDeadline, referenceDate))

	// ── Year-to-date window ─────────────────────────────────────────────────
	ytdStart, ytdEnd := ytdPeriod(referenceDate)
	ytdDays := timefy.New(ytdEnd).DurationInDays(ytdStart) + 1
	fmt.Printf("\n  YTD period   : %s → %s  (%d days)\n",
		ytdStart.Format("2006-01-02"), ytdEnd.Format("2006-01-02 15:04:05"),
		ytdDays)

	// ── Quarter-over-quarter comparison window ──────────────────────────────
	curr, prior := qoqComparison(referenceDate)
	fmt.Println("\n  Quarter-over-quarter comparison:")
	fmt.Printf("  Current  : %-10s  %s → %s\n", curr.Label,
		curr.Start.Format("2006-01-02"), curr.End.Format("2006-01-02"))
	fmt.Printf("  Prior YoY: %-10s  %s → %s\n", prior.Label,
		prior.Start.Format("2006-01-02"), prior.End.Format("2006-01-02"))

	// ── All 4 calendar quarters of the current year ─────────────────────────
	fmt.Printf("\n  All calendar quarters of %d:\n", referenceDate.Year())
	fmt.Println("  ─────────────────────────────────────────────────────────")
	for q := 1; q <= 4; q++ {
		// Pick the first month of each quarter.
		month := time.Month((q-1)*3 + 1)
		rep := time.Date(referenceDate.Year(), month, 1, 0, 0, 0, 0, time.UTC)
		qp := calendarQuarterPeriod(rep)
		days := timefy.New(qp.End).DurationInDays(qp.Start) + 1
		fmt.Printf("  %-10s  %s → %s  (%d days)  filing by %s\n",
			qp.Label,
			qp.Start.Format("Jan 02"),
			qp.End.Format("Jan 02"),
			days,
			qp.FilingDeadline.Format("2006-01-02"),
		)
	}

	// ── All 4 fiscal quarters of FY2025 ────────────────────────────────────
	fmt.Printf("\n  All fiscal quarters of FY2025 (Jul 2024 – Jun 2025):\n")
	fmt.Println("  ─────────────────────────────────────────────────────────")
	fiscalRefDates := []time.Time{
		time.Date(2024, time.July, 15, 0, 0, 0, 0, time.UTC),    // FQ1
		time.Date(2024, time.October, 15, 0, 0, 0, 0, time.UTC), // FQ2
		time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC), // FQ3
		time.Date(2025, time.April, 15, 0, 0, 0, 0, time.UTC),   // FQ4
	}
	for _, ref := range fiscalRefDates {
		fqp := fiscalQuarterPeriod(ref)
		days := timefy.New(fqp.End).DurationInDays(fqp.Start) + 1
		fmt.Printf("  %-12s  %s → %s  (%d days)  filing by %s\n",
			fqp.Label,
			fqp.Start.Format("2006-01-02"),
			fqp.End.Format("2006-01-02"),
			days,
			fqp.FilingDeadline.Format("2006-01-02"),
		)
	}

	// ── Weekdays in Q1 2025 (for trading-day count) ──────────────────────────
	q1Start := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	q1End := time.Date(2025, time.March, 31, 0, 0, 0, 0, time.UTC)
	weekdays := timefy.GetWeekdaysInRange(q1Start, q1End)
	fmt.Printf("\n  Trading days in Q1 2025 (Mon–Fri): %d\n", len(weekdays))

	// ── Half-year reporting window (H1 / H2) ─────────────────────────────────
	h1End := timefy.New(referenceDate).BeginningOfHalf()
	h1Start := timefy.New(h1End).BeginningOfYear()
	fmt.Printf("\n  H1 %d boundaries  : %s → %s\n",
		referenceDate.Year(),
		h1Start.Format("2006-01-02"),
		h1End.AddDate(0, 6, 0).Add(-time.Nanosecond).Format("2006-01-02"))

	// ── Human-readable time until year end ───────────────────────────────────
	yearEnd := timefy.New(referenceDate).EndOfYear()
	fmt.Printf("\n  Year-end %d countdown : %s\n",
		referenceDate.Year(),
		timefy.New(yearEnd).TimeUntil())

	// ── Leap-year check using IsLeapYear ─────────────────────────────────────
	fmt.Println("\n  Leap-year check for recent years:")
	for _, yr := range []int{2020, 2021, 2022, 2023, 2024, 2025} {
		leapFlag := timefy.IsLeapYear(yr)
		t := time.Date(yr, time.February, 1, 0, 0, 0, 0, time.UTC)
		febEnd := timefy.New(t).EndOfMonth()
		fmt.Printf("  %d  leap=%-5v  end-of-Feb=%s\n",
			yr, leapFlag, febEnd.Format("2006-01-02"))
	}
}
