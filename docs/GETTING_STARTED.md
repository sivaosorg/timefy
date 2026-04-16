# timefy — Comprehensive Technical Documentation

> **Senior-level reference for Go engineers using the `timefy` library.**
> This documentation teaches the library from first principles through advanced
> enterprise patterns, with every concept backed by complete, runnable Go code.

---

## Table of Contents

1. [Introduction & Design Philosophy](#1-introduction--design-philosophy)
2. [Architecture Overview](#2-architecture-overview)
3. [Installation](#3-installation)
4. [Core Type System](#4-core-type-system)
   - 4.1 [TimeRFC — Short Time Layouts](#41-timerfc--short-time-layouts)
   - 4.2 [TimeFormatRFC — Full Date-Time Layouts](#42-timeformatrfc--full-date-time-layouts)
   - 4.3 [ZoneRFC — IANA Timezone Identifiers](#43-zonerfc--iana-timezone-identifiers)
   - 4.4 [Rule — Configuration Object](#44-rule--configuration-object)
   - 4.5 [Timex — The Enriched Time Value](#45-timex--the-enriched-time-value)
5. [Getting Started: Wrapping Time Values](#5-getting-started-wrapping-time-values)
6. [Time Boundaries](#6-time-boundaries)
   - 6.1 [Minute and Hour Boundaries](#61-minute-and-hour-boundaries)
   - 6.2 [Day Boundaries](#62-day-boundaries)
   - 6.3 [Week Boundaries](#63-week-boundaries)
   - 6.4 [Month Boundaries](#64-month-boundaries)
   - 6.5 [Quarter Boundaries](#65-quarter-boundaries)
   - 6.6 [Half-Year Boundaries](#66-half-year-boundaries)
   - 6.7 [Year Boundaries](#67-year-boundaries)
   - 6.8 [Global Convenience Functions](#68-global-convenience-functions)
7. [Day Operations](#7-day-operations)
8. [Time Arithmetic](#8-time-arithmetic)
9. [Parsing Strings into Time](#9-parsing-strings-into-time)
   - 9.1 [The Parse Contract](#91-the-parse-contract)
   - 9.2 [Parse — Safe Variant](#92-parse--safe-variant)
   - 9.3 [MustParse — Panic-on-Failure Variant](#93-mustparse--panic-on-failure-variant)
   - 9.4 [ParseInLocation and MustParseInLocation](#94-parseinlocation-and-mustparseinlocation)
   - 9.5 [Extending the Format Registry](#95-extending-the-format-registry)
10. [Formatting Time](#10-formatting-time)
    - 10.1 [Go's Reference Time](#101-gos-reference-time)
    - 10.2 [FormatRFC and FormatRFCshort](#102-formatrfc-and-formatrfcshort)
    - 10.3 [FormatTimex — Component Extraction](#103-formattimex--component-extraction)
11. [Timezone Operations](#11-timezone-operations)
    - 11.1 [SetTimezone vs AdjustTimezone](#111-settimezone-vs-adjusttimezone)
    - 11.2 [Predefined ZoneRFC Constants](#112-predefined-zonerfc-constants)
    - 11.3 [Rule-based Timezone Configuration](#113-rule-based-timezone-configuration)
12. [Weekday Navigation](#12-weekday-navigation)
13. [Human-Readable Relative Time](#13-human-readable-relative-time)
    - 13.1 [TimeAgo](#131-timeago)
    - 13.2 [TimeUntil](#132-timeuntil)
14. [Duration Calculations](#14-duration-calculations)
15. [Time Comparisons](#15-time-comparisons)
    - 15.1 [Ordinal Comparisons](#151-ordinal-comparisons)
    - 15.2 [Range Membership](#152-range-membership)
    - 15.3 [Calendrical Identity](#153-calendrical-identity)
    - 15.4 [Temporal Position](#154-temporal-position)
    - 15.5 [Day-Type Checks](#155-day-type-checks)
16. [Between (String-based Range Check)](#16-between-string-based-range-check)
17. [Elapsed Time Measurements](#17-elapsed-time-measurements)
18. [Leap Year Utilities](#18-leap-year-utilities)
19. [Weekdays in Range](#19-weekdays-in-range)
20. [Tolerance Window](#20-tolerance-window)
21. [Quarter Number](#21-quarter-number)
22. [Rule Configuration — Full API](#22-rule-configuration--full-api)
23. [Advanced Patterns](#23-advanced-patterns)
    - 23.1 [Fluent Builder Chains](#231-fluent-builder-chains)
    - 23.2 [Multi-tenant Timezone Isolation](#232-multi-tenant-timezone-isolation)
    - 23.3 [Custom Format Registries](#233-custom-format-registries)
    - 23.4 [Safe vs Panic Patterns](#234-safe-vs-panic-patterns)
24. [Enterprise Use Cases](#24-enterprise-use-cases)
25. [API Quick-Reference](#25-api-quick-reference)
26. [Format Constant Reference](#26-format-constant-reference)
27. [Timezone Constant Reference](#27-timezone-constant-reference)

---

## 1. Introduction & Design Philosophy

`timefy` is a zero-dependency Go library that extends the standard `time`
package with a rich, ergonomic API for time management tasks that recur
constantly in production software.

### Why `timefy` exists

The standard library's `time.Time` is precise and correct, but deliberately
minimal.  Everyday tasks require boilerplate that is easy to get wrong:

| Task | Standard Library | timefy |
|------|-----------------|--------|
| Start of current month | `time.Date(y, m, 1, 0, 0, 0, 0, loc)` | `With(now).BeginningOfMonth()` |
| Last day of quarter | `complex 3-step calculation` | `With(now).EndOfQuarter()` |
| Parse any of 40+ formats | `try each format manually` | `Parse("2025-03-15")` |
| "5 minutes ago" | `manual switch-case` | `New(t).TimeAgo()` |
| Is this a weekday? | `t.Weekday() != Sat && != Sun` | `New(t).IsWeekday()` |

### Design decisions

1. **Immutability**: all boundary methods return new `time.Time` values;
   the wrapped `Timex` is never mutated.
2. **Nanosecond precision at boundaries**: `EndOfDay` returns
   `23:59:59.999999999`, not `23:59:59.000000000`.  This is critical for
   inclusive range queries in databases.
3. **Zero external dependencies**: the module's only import is the standard
   library.
4. **Fluent configuration via Rule**: timezone, week-start, and format registry
   can be configured per-call via a `Rule` builder without touching globals.
5. **Two-tier API**: global convenience functions (`BeginningOfDay()`) operate
   on `time.Now()`; method-based API (`With(t).BeginningOfDay()`) operates on
   any arbitrary time.

---

## 2. Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                       Public API Surface                     │
│                                                             │
│  Global functions      Rule (config)       Timex methods   │
│  ─────────────────     ──────────────      ─────────────── │
│  BeginningOfDay()      NewRule()           .BeginningOfDay()│
│  Parse("2025-03")      .WithLocation()     .EndOfMonth()   │
│  AddHour(t, 3)         .With(t) → Timex    .TimeAgo()      │
│  IsLeapYear(2024)      .Parse("…")         .IsBetween()    │
└───────────────────────────┬─────────────────────────────────┘
                            │
                   ┌────────▼────────┐
                   │   type Timex    │
                   │  ┌───────────┐  │
                   │  │ time.Time │  │  ← embedded; inherits all
                   │  └───────────┘  │    standard methods
                   │  ┌───────────┐  │
                   │  │  *Rule    │  │  ← configuration context
                   │  └───────────┘  │
                   └─────────────────┘
                            │
              ┌─────────────▼──────────────┐
              │          type Rule          │
              │  weekStartDay time.Weekday  │
              │  timeLocation *time.Location│
              │  timeFormats  []string      │
              └─────────────────────────────┘
```

The `Timex` struct embeds `time.Time` (inheriting all 40+ standard methods)
and carries a `*Rule` pointer that governs:

- Which day the week starts on (`weekStartDay`)
- Which timezone to use when evaluating boundaries (`timeLocation`)
- Which format strings to try when parsing (`timeFormats`)

---

## 3. Installation

```bash
# Specific version
go get github.com/sivaosorg/timefy@v0.0.3

# Latest
go get github.com/sivaosorg/timefy@latest
```

```go
import "github.com/sivaosorg/timefy"
```

**Minimum Go version:** 1.23.1

---

## 4. Core Type System

### 4.1 TimeRFC — Short Time Layouts

```go
type TimeRFC string
```

`TimeRFC` is a typed string alias for short (time-only) format patterns.
Using a distinct type prevents passing a full date-time layout where a
short layout is expected, catching mistakes at compile time.

```go
const (
    TimeRFC01T150405 TimeRFC = "15:04:05"   // HH:MM:SS (colon-separated)
    TimeRFC02D150405 TimeRFC = "15.04.05"   // HH.MM.SS (dot-separated)
    TimeRFC03H150405 TimeRFC = "15-04-05"   // HH-MM-SS (dash-separated)
    TimeRFC04M150405 TimeRFC = "15/04/05"   // HH/MM/SS (slash-separated)
    TimeRFC05S150405 TimeRFC = "15 04 05"   // HH MM SS (space-separated)
    TimeRFC06T150405 TimeRFC = "150405"     // HHMMSS (no separator)
    TimeRFC07D150405 TimeRFC = "15:04"      // HH:MM
    TimeRFC08H150405 TimeRFC = "15.04"      // HH.MM
    TimeRFC09M150405 TimeRFC = "15-04"      // HH-MM
    TimeRFC10S150405 TimeRFC = "15/04"      // HH/MM
    TimeRFC11T150405 TimeRFC = "1504"       // HHMM
)
```

### 4.2 TimeFormatRFC — Full Date-Time Layouts

```go
type TimeFormatRFC string
```

`TimeFormatRFC` covers full date-time formats (date + optional time).
Selected constants:

```go
const (
    TimeFormat20060102T150405999999   TimeFormatRFC = "2006-01-02T15:04:05.999999"
    TimeFormat20060102T150405         TimeFormatRFC = "2006-01-02T15:04:05"
    TimeFormat20060102150405          TimeFormatRFC = "2006-01-02 15:04:05"
    TimeFormat20060102                TimeFormatRFC = "2006-01-02"
    TimeFormat200601                  TimeFormatRFC = "2006-01"
    TimeFormat20060102150405999999    TimeFormatRFC = "2006-01-02 15:04:05.999999"
    TimeFormat20060102150405Z0700     TimeFormatRFC = "2006-01-02 15:04:05 -07:00"
    TimeFormat20060102150405Z0700RFC3339 TimeFormatRFC = "2006-01-02T15:04:05-07:00"
)
```

### 4.3 ZoneRFC — IANA Timezone Identifiers

```go
type ZoneRFC string
```

`ZoneRFC` wraps an IANA timezone name string.  The type prevents accidentally
passing a format string where a timezone name is expected.

```go
const (
    DefaultTimezoneVietnam    ZoneRFC = "Asia/Ho_Chi_Minh"
    DefaultTimezoneNewYork    ZoneRFC = "America/New_York"
    DefaultTimezoneLondon     ZoneRFC = "Europe/London"
    DefaultTimezoneTokyo      ZoneRFC = "Asia/Tokyo"
    DefaultTimezoneSydney     ZoneRFC = "Australia/Sydney"
    DefaultTimezoneBeijing    ZoneRFC = "Asia/Shanghai"
    DefaultTimezoneSingapore  ZoneRFC = "Asia/Singapore"
    // ... 30+ more constants
)
```

### 4.4 Rule — Configuration Object

```go
type Rule struct {
    weekStartDay time.Weekday    // default: time.Sunday
    timeLocation *time.Location  // default: nil (uses time value's location)
    timeFormats  []string        // default: 40+ built-in patterns
}
```

`Rule` is constructed with `NewRule()` and configured through a fluent
builder pattern.  It is then used to create `Timex` values via `rule.With(t)`.

```go
rule := timefy.NewRule().
    WithWeekStartDay(time.Monday).
    WithLocationRFC(timefy.DefaultTimezoneVietnam)

timex := rule.With(time.Now())
```

### 4.5 Timex — The Enriched Time Value

```go
type Timex struct {
    time.Time        // embedded: inherits all standard time.Time methods
    *Rule            // pointer to configuration context
}
```

Because `time.Time` is embedded, every `time.Time` method is available
directly on `*Timex`:

```go
tx := timefy.With(time.Now())
fmt.Println(tx.Year())            // standard time.Time.Year()
fmt.Println(tx.BeginningOfMonth()) // timefy extension
```

---

## 5. Getting Started: Wrapping Time Values

There are two equivalent ways to create a `Timex` with default configuration:

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    // With() and New() are identical; both wrap a time.Time into a *Timex.
    // The choice between them is stylistic.
    now := time.Now()

    tx1 := timefy.With(now) // preferred when chaining from a variable
    tx2 := timefy.New(now)  // preferred when reading as "new Timex for now"

    // Both produce a *Timex that embeds the given time.Time and
    // a default Rule (Sunday week-start, system formats, nil location).
    fmt.Println(tx1.BeginningOfDay())
    fmt.Println(tx2.BeginningOfDay())

    // To use a specific time value (not now):
    t := time.Date(2025, 1, 15, 17, 51, 49, 123456789, time.UTC)
    tx := timefy.New(t)
    fmt.Println(tx.EndOfMonth()) // 2025-01-31 23:59:59.999999999 +0000 UTC
}
```

**Why not operate directly on `time.Time`?**

`time.Time` has no concept of week-start day or custom format registries.
`Timex` carries this context throughout a chain of operations without
requiring the caller to pass configuration on every call.

---

## 6. Time Boundaries

Time boundaries are one of the most error-prone areas of date arithmetic.
Common pitfalls include:

- Off-by-one errors at month/year boundaries (e.g., "March 32")
- Forgetting that `EndOfMonth` varies (28/29/30/31 days)
- Forgetting nanosecond precision for inclusive `BETWEEN` queries
- DST transitions changing the nominal clock time of a boundary

`timefy` handles all of these correctly.

### 6.1 Minute and Hour Boundaries

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    t := time.Date(2025, 8, 15, 13, 45, 30, 500000000, time.UTC)
    tx := timefy.New(t)

    // BeginningOfMinute: truncates to the nearest minute.
    // Internally: t.Truncate(time.Minute)
    // Result: 2025-08-15 13:45:00 UTC
    fmt.Println(tx.BeginningOfMinute())

    // EndOfMinute: BeginningOfMinute + 1 minute - 1 nanosecond.
    // The "- 1 nanosecond" ensures the boundary is non-overlapping:
    // a time at 13:46:00.000 is NOT included in the 13:45 minute.
    // Result: 2025-08-15 13:45:59.999999999 UTC
    fmt.Println(tx.EndOfMinute())

    // BeginningOfHour: resets minute, second, nanosecond to zero.
    // Result: 2025-08-15 13:00:00 UTC
    fmt.Println(tx.BeginningOfHour())

    // EndOfHour: BeginningOfHour + 1 hour - 1 nanosecond.
    // Result: 2025-08-15 13:59:59.999999999 UTC
    fmt.Println(tx.EndOfHour())
}
```

### 6.2 Day Boundaries

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    t := time.Date(2025, 3, 15, 17, 30, 0, 0, time.UTC)
    tx := timefy.New(t)

    // BeginningOfDay: midnight 00:00:00.000000000
    fmt.Println(tx.BeginningOfDay()) // 2025-03-15 00:00:00 +0000 UTC

    // EndOfDay: 23:59:59.999999999 — the last representable nanosecond
    // of the day.  This is critical for SQL BETWEEN queries:
    //   WHERE created_at BETWEEN '2025-03-15 00:00:00' AND '2025-03-15 23:59:59.999999999'
    // will match all records on March 15 including those at 23:59:59.
    fmt.Println(tx.EndOfDay()) // 2025-03-15 23:59:59.999999999 +0000 UTC

    // Package-level equivalents using a concrete time.Time (not Timex):
    fmt.Println(timefy.BeginOfDay(t))  // 2025-03-15 00:00:00 +0000 UTC
    fmt.Println(timefy.FEndOfDay(t))   // 2025-03-15 23:59:59 +0000 UTC
    // Note: FEndOfDay truncates at seconds (no nanoseconds); prefer
    // New(t).EndOfDay() when nanosecond precision is required.

    // Previous day boundaries:
    fmt.Println(timefy.PrevBeginOfDay(t, 1)) // yesterday start
    fmt.Println(timefy.PrevEndOfDay(t, 1))   // yesterday end
    fmt.Println(timefy.PrevBeginOfDay(t, 7)) // one week ago start
}
```

### 6.3 Week Boundaries

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    // Wednesday, March 19, 2025
    t := time.Date(2025, 3, 19, 12, 0, 0, 0, time.UTC)

    // Default week-start is Sunday.
    tx := timefy.New(t)
    fmt.Println(tx.BeginningOfWeek()) // 2025-03-16 00:00:00 (Sunday)
    fmt.Println(tx.EndOfWeek())       // 2025-03-22 23:59:59.999999999 (Saturday)

    // Configure Monday as week-start for ISO 8601 calendars
    // (European business standard).
    rule := timefy.NewRule().WithWeekStartDay(time.Monday)
    txISO := rule.With(t)
    fmt.Println(txISO.BeginningOfWeek()) // 2025-03-17 00:00:00 (Monday)
    fmt.Println(txISO.EndOfWeek())       // 2025-03-23 23:59:59.999999999 (Sunday)
}
```

**Why week-start matters:**
The ISO 8601 standard defines the week as Monday–Sunday.  Many European
companies, regulatory bodies, and international standards use this calendar.
The default Sunday–Saturday week is the US convention.  `timefy` allows
both via `WithWeekStartDay`.

### 6.4 Month Boundaries

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    t := time.Date(2025, 3, 15, 12, 0, 0, 0, time.UTC)
    tx := timefy.New(t)

    fmt.Println(tx.BeginningOfMonth()) // 2025-03-01 00:00:00 +0000 UTC
    fmt.Println(tx.EndOfMonth())       // 2025-03-31 23:59:59.999999999 +0000 UTC

    // EndOfMonth correctly handles variable-length months:
    months := []time.Time{
        time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC), // January (31 days)
        time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC), // February 2025 (28 days)
        time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC), // February 2024 (29 days, leap)
        time.Date(2025, 4, 15, 0, 0, 0, 0, time.UTC), // April (30 days)
    }
    for _, m := range months {
        end := timefy.New(m).EndOfMonth()
        fmt.Printf("  %s → last day: %s (day %d)\n",
            m.Format("Jan 2006"), end.Format("2006-01-02"), end.Day())
    }
    // Output:
    //   Jan 2025 → last day: 2025-01-31 (day 31)
    //   Feb 2025 → last day: 2025-02-28 (day 28)
    //   Feb 2024 → last day: 2024-02-29 (day 29)
    //   Apr 2025 → last day: 2025-04-30 (day 30)
}
```

**Implementation insight:** `EndOfMonth` is computed as:
```
BeginningOfMonth().AddDate(0, 1, 0).Add(-time.Nanosecond)
```
By adding 1 month to the first day and subtracting 1 nanosecond, the
month-length calculation is delegated entirely to Go's `AddDate`, which
handles all month-length and leap-year arithmetic correctly.

### 6.5 Quarter Boundaries

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    quarters := []time.Time{
        time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC), // Q1 (Jan-Mar)
        time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC), // Q2 (Apr-Jun)
        time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC),  // Q3 (Jul-Sep)
        time.Date(2025, 11, 20, 0, 0, 0, 0, time.UTC), // Q4 (Oct-Dec)
    }

    for _, t := range quarters {
        tx := timefy.New(t)
        fmt.Printf("  %s → Q%d: %s → %s\n",
            t.Format("Jan 02"),
            tx.Quarter(),
            tx.BeginningOfQuarter().Format("2006-01-02"),
            tx.EndOfQuarter().Format("2006-01-02"),
        )
    }
    // Output:
    //   Feb 15 → Q1: 2025-01-01 → 2025-03-31
    //   May 10 → Q2: 2025-04-01 → 2025-06-30
    //   Sep 01 → Q3: 2025-07-01 → 2025-09-30
    //   Nov 20 → Q4: 2025-10-01 → 2025-12-31
}
```

**Implementation insight:** `BeginningOfQuarter` computes:
```
offset = (currentMonth - 1) % 3   // months into the current quarter
start  = BeginningOfMonth().AddDate(0, -offset, 0)
```
This yields Q1 for month 1 (offset=0), month 2 (offset=1), month 3 (offset=2).

### 6.6 Half-Year Boundaries

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    t := time.Date(2025, 9, 15, 0, 0, 0, 0, time.UTC) // H2 (Jul-Dec)
    tx := timefy.New(t)

    fmt.Println(tx.BeginningOfHalf()) // 2025-07-01 00:00:00
    fmt.Println(tx.EndOfHalf())       // 2025-12-31 23:59:59.999999999

    t2 := time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC) // H1 (Jan-Jun)
    tx2 := timefy.New(t2)
    fmt.Println(tx2.BeginningOfHalf()) // 2025-01-01 00:00:00
    fmt.Println(tx2.EndOfHalf())       // 2025-06-30 23:59:59.999999999
}
```

Half-year boundaries are used in financial reporting (H1/H2 earnings) and
performance review cycles (mid-year reviews).

### 6.7 Year Boundaries

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    t := time.Date(2025, 7, 4, 10, 0, 0, 0, time.UTC)
    tx := timefy.New(t)

    fmt.Println(tx.BeginningOfYear()) // 2025-01-01 00:00:00 +0000 UTC
    fmt.Println(tx.EndOfYear())       // 2025-12-31 23:59:59.999999999 +0000 UTC
}
```

### 6.8 Global Convenience Functions

All `BeginningOf*` and `EndOf*` methods have package-level counterparts that
implicitly use `time.Now()`:

```go
package main

import (
    "fmt"

    "github.com/sivaosorg/timefy"
)

func main() {
    // These are equivalent to timefy.With(time.Now()).BeginningOfDay(), etc.
    fmt.Println(timefy.BeginningOfMinute())
    fmt.Println(timefy.BeginningOfHour())
    fmt.Println(timefy.BeginningOfDay())
    fmt.Println(timefy.BeginningOfWeek())
    fmt.Println(timefy.BeginningOfMonth())
    fmt.Println(timefy.BeginningOfQuarter())
    fmt.Println(timefy.BeginningOfYear())

    fmt.Println(timefy.EndOfMinute())
    fmt.Println(timefy.EndOfHour())
    fmt.Println(timefy.EndOfDay())
    fmt.Println(timefy.EndOfWeek())
    fmt.Println(timefy.EndOfMonth())
    fmt.Println(timefy.EndOfQuarter())
    fmt.Println(timefy.EndOfYear())
}
```

---

## 7. Day Operations

Package-level functions for concrete `time.Time` day manipulation (distinct
from the `Timex` boundary methods in § 6):

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    t := time.Date(2025, 3, 15, 17, 30, 0, 0, time.UTC)

    // BeginOfDay: equivalent to Timex.BeginningOfDay but takes a time.Time.
    fmt.Println(timefy.BeginOfDay(t))  // 2025-03-15 00:00:00 +0000 UTC

    // FEndOfDay: end of day at second precision (23:59:59, no nanoseconds).
    // Use Timex.EndOfDay() when nanosecond precision is required.
    fmt.Println(timefy.FEndOfDay(t))   // 2025-03-15 23:59:59 +0000 UTC

    // PrevBeginOfDay(t, n): start of the day n days before t.
    fmt.Println(timefy.PrevBeginOfDay(t, 1)) // 2025-03-14 00:00:00 (yesterday)
    fmt.Println(timefy.PrevBeginOfDay(t, 7)) // 2025-03-08 00:00:00 (7 days ago)
    fmt.Println(timefy.PrevBeginOfDay(t, 30)) // 2025-02-13 00:00:00 (30 days ago)

    // PrevEndOfDay(t, n): end of the day n days before t.
    fmt.Println(timefy.PrevEndOfDay(t, 1)) // 2025-03-14 23:59:59 (end of yesterday)
}
```

---

## 8. Time Arithmetic

Add or subtract precise time units from any `time.Time`:

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    t := time.Date(2025, 3, 15, 12, 0, 0, 0, time.UTC)

    // AddSecond(t, n): add n seconds (negative = subtract).
    // If n == 0, returns t unchanged (optimization).
    fmt.Println(timefy.AddSecond(t, 30))   // 2025-03-15 12:00:30
    fmt.Println(timefy.AddSecond(t, -90))  // 2025-03-15 11:58:30
    fmt.Println(timefy.AddSecond(t, 0))    // 2025-03-15 12:00:00 (unchanged)

    // AddMinute(t, n): add n minutes.
    fmt.Println(timefy.AddMinute(t, 90))   // 2025-03-15 13:30:00
    fmt.Println(timefy.AddMinute(t, -45))  // 2025-03-15 11:15:00

    // AddHour(t, n): add n hours.
    fmt.Println(timefy.AddHour(t, 12))     // 2025-03-16 00:00:00 (crosses midnight)
    fmt.Println(timefy.AddHour(t, -24))    // 2025-03-14 12:00:00

    // AddDay(t, n): add n × 24 hours.
    // Note: uses time.Add(24h * n), NOT AddDate(0, 0, n).
    // For DST-sensitive applications, prefer time.AddDate(0, 0, n) when
    // calendar-day semantics are required across DST transitions.
    fmt.Println(timefy.AddDay(t, 15))      // 2025-03-30 12:00:00
    fmt.Println(timefy.AddDay(t, -15))     // 2025-02-28 12:00:00
}
```

---

## 9. Parsing Strings into Time

### 9.1 The Parse Contract

`timefy` maintains a global `TimeFormats` slice containing 40+ patterns.
When parsing, the library tries each pattern in order and returns the
first successful match.  Missing components are filled in from the reference
time (typically `time.Now()`).

```go
// TimeFormats (excerpt):
var TimeFormats = []string{
    "2006",               // Year only
    "2006-1",             // Year-Month
    "2006-1-2",           // Year-Month-Day
    "2006-1-2 15:4:5",    // Full date-time
    "15:4:5",             // Time only
    "15:4",               // Hour:Minute
    "15",                 // Hour only
    time.RFC3339,
    time.RFC3339Nano,
    // ... 30+ more
}
```

### 9.2 Parse — Safe Variant

```go
package main

import (
    "fmt"

    "github.com/sivaosorg/timefy"
)

func main() {
    // Full date-time:
    t, err := timefy.Parse("2025-03-15 14:30:00")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(t) // 2025-03-15 14:30:00 +local

    // Date only — time components default to 00:00:00:
    t, _ = timefy.Parse("2025-03-15")
    fmt.Println(t) // 2025-03-15 00:00:00

    // Year-Month — day defaults to 01:
    t, _ = timefy.Parse("2025-03")
    fmt.Println(t) // 2025-03-01 00:00:00

    // Year only — defaults to January 1st:
    t, _ = timefy.Parse("2025")
    fmt.Println(t) // 2025-01-01 00:00:00

    // Time only — date is taken from today:
    t, _ = timefy.Parse("14:30")
    fmt.Println(t) // <today> 14:30:00

    // Hour only:
    t, _ = timefy.Parse("14")
    fmt.Println(t) // <today> 14:00:00

    // Month-Day (no year) — year taken from today:
    t, _ = timefy.Parse("11-15")
    fmt.Println(t) // <current year>-11-15 00:00:00

    // Error case:
    _, err = timefy.Parse("not-a-time")
    fmt.Println(err) // can't parse string as time: not-a-time
}
```

### 9.3 MustParse — Panic-on-Failure Variant

```go
package main

import (
    "fmt"

    "github.com/sivaosorg/timefy"
)

func main() {
    // MustParse panics if parsing fails.
    // Use only when the input is guaranteed valid (config files, constants).
    t := timefy.MustParse("2025-03-15")
    fmt.Println(t) // 2025-03-15 00:00:00

    // This will PANIC — do not use with user-supplied input:
    // timefy.MustParse("99:99") // panic: can't parse string as time: 99:99
}
```

### 9.4 ParseInLocation and MustParseInLocation

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    // ParseInLocation interprets the string in the specified timezone.
    // The parsed time's Location is set to loc.
    t, err := timefy.ParseInLocation(time.UTC, "2025-03-15 14:30:00")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(t) // 2025-03-15 14:30:00 +0000 UTC

    // Compare with Parse() which uses the system local timezone:
    t2, _ := timefy.Parse("2025-03-15 14:30:00")
    fmt.Println(t2) // 2025-03-15 14:30:00 +0700 +07  (if local = +7)

    // Load any IANA timezone:
    loc, _ := time.LoadLocation("America/New_York")
    t3 := timefy.MustParseInLocation(loc, "2025-03-15 14:30:00")
    fmt.Println(t3) // 2025-03-15 14:30:00 -0400 EDT
}
```

### 9.5 Extending the Format Registry

The global `TimeFormats` variable is mutable, allowing custom formats to be
registered at program start:

```go
package main

import (
    "fmt"

    "github.com/sivaosorg/timefy"
)

func init() {
    // Add European date format "02 Jan 2006 15:04" once at startup.
    timefy.TimeFormats = append(timefy.TimeFormats, "02 Jan 2006 15:04")
}

func main() {
    t, err := timefy.Parse("15 Mar 2025 14:30")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(t) // 2025-03-15 14:30:00
}
```

Alternatively, use a `Rule` to scope custom formats to a single call site
(see § 22).

---

## 10. Formatting Time

### 10.1 Go's Reference Time

Go's time formatting works by example, not by `%Y-%m-%d` directives.
The reference timestamp is:

```
Mon Jan 2 15:04:05 MST 2006
```

which encoded numerically is `01/02 03:04:05PM '06 -0700`.  Every format
string in `timefy` (and in Go) uses this reference.

```
Year:  2006          Month:  01          Day:    02
Hour:  15 (24h)      Minute: 04          Second: 05
```

### 10.2 FormatRFC and FormatRFCshort

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    t := time.Date(2025, 8, 15, 13, 45, 30, 123456789, time.UTC)
    tx := timefy.New(t)

    // FormatRFC: format using a TimeFormatRFC constant.
    fmt.Println(tx.FormatRFC(timefy.TimeFormat20060102T150405))
    // → 2025-08-15T13:45:30

    fmt.Println(tx.FormatRFC(timefy.TimeFormat20060102T150405999999))
    // → 2025-08-15T13:45:30.123456

    fmt.Println(tx.FormatRFC(timefy.TimeFormat20060102))
    // → 2025-08-15

    // FormatRFCshort: format using a TimeRFC constant (time only).
    fmt.Println(tx.FormatRFCshort(timefy.TimeRFC01T150405))
    // → 13:45:30

    fmt.Println(tx.FormatRFCshort(timefy.TimeRFC07D150405))
    // → 13:45

    // DefaultFormatRFC: "2006-01-02 15:04:05.999999" (microsecond precision).
    fmt.Println(tx.DefaultFormatRFC())
    // → 2025-08-15 13:45:30.123456

    // Package-level variants (format time.Now()):
    fmt.Println(timefy.FormatRFC(timefy.TimeFormat20060102T150405))
    fmt.Println(timefy.FormatRFCshort(timefy.TimeRFC01T150405))
    fmt.Println(timefy.DefaultFormatRFC())

    // Package-level variants for a specific time.Time:
    fmt.Println(timefy.FFormatRFC(t, timefy.TimeFormat20060102T150405))
    fmt.Println(timefy.FFormatRFCshort(t, timefy.TimeRFC01T150405))
    fmt.Println(timefy.FDefaultFormatRFC(t))
}
```

### 10.3 FormatTimex — Component Extraction

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    t := time.Date(2025, 3, 15, 14, 30, 45, 123456789, time.UTC)

    // FormatTimex returns []int{nanosecond, second, minute, hour, day, month, year}
    c := timefy.FormatTimex(t)
    // c = [123456789, 45, 30, 14, 15, 3, 2025]
    //      ns         sec min  hr  day mo  year

    fmt.Printf("Year:   %d\n", c[6])  // 2025
    fmt.Printf("Month:  %d\n", c[5])  // 3
    fmt.Printf("Day:    %d\n", c[4])  // 15
    fmt.Printf("Hour:   %d\n", c[3])  // 14
    fmt.Printf("Minute: %d\n", c[2])  // 30
    fmt.Printf("Second: %d\n", c[1])  // 45
    fmt.Printf("Nano:   %d\n", c[0])  // 123456789

    // FormatTimex is used internally by Timex.Parse to reconstruct a
    // time.Time from partially-parsed components.  It is also useful
    // in structured logging where individual components are separate fields.
}
```

---

## 11. Timezone Operations

### 11.1 SetTimezone vs AdjustTimezone

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    utc := time.Date(2025, 3, 15, 12, 0, 0, 0, time.UTC)

    // SetTimezone: converts the time to a new timezone.
    // Returns an error if the timezone name is invalid.
    nyTime, err := timefy.SetTimezone(utc, "America/New_York")
    if err != nil {
        fmt.Println("Error:", err) // e.g., "unknown time zone Foo/Bar"
        return
    }
    fmt.Println(nyTime)  // 2025-03-15 08:00:00 -0400 EDT (UTC-4 in DST)

    // AdjustTimezone: same as SetTimezone but returns the original time
    // unchanged if the timezone name is invalid (no error return).
    // Use this when timezone validity is guaranteed.
    tokyoTime := timefy.AdjustTimezone(utc, "Asia/Tokyo")
    fmt.Println(tokyoTime) // 2025-03-15 21:00:00 +0900 JST

    // Invalid timezone: original time returned silently.
    invalid := timefy.AdjustTimezone(utc, "Invalid/Zone")
    fmt.Println(invalid) // 2025-03-15 12:00:00 +0000 UTC (unchanged)

    // The two functions differ in error handling strategy:
    //   SetTimezone  → explicit error → use in validation paths
    //   AdjustTimezone → silent fallback → use in rendering paths
}
```

### 11.2 Predefined ZoneRFC Constants

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    utc := time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC)

    zones := []timefy.ZoneRFC{
        timefy.DefaultTimezoneVietnam,    // UTC+7
        timefy.DefaultTimezoneNewYork,    // UTC-5 / UTC-4 (DST)
        timefy.DefaultTimezoneLondon,     // UTC+0 / UTC+1 (BST)
        timefy.DefaultTimezoneTokyo,      // UTC+9
        timefy.DefaultTimezoneSydney,     // UTC+10 / UTC+11 (AEDT)
        timefy.DefaultTimezoneDelhi,      // UTC+5:30
        timefy.DefaultTimezoneBeijing,    // UTC+8
        timefy.DefaultTimezoneSingapore,  // UTC+8
        timefy.DefaultTimezoneDubai,      // UTC+4
        timefy.DefaultTimezoneMoscow,     // UTC+3
    }

    for _, z := range zones {
        local := timefy.AdjustTimezone(utc, string(z))
        fmt.Printf("  %-30s → %s\n", z, local.Format("15:04 MST"))
    }
}
```

### 11.3 Rule-based Timezone Configuration

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    // Configure a Rule pinned to Vietnam time.
    rule := timefy.NewRule().WithLocationRFC(timefy.DefaultTimezoneVietnam)

    // rule.Parse evaluates the string in Vietnam timezone context.
    t, _ := rule.Parse("2025-11-12 22:14:01")
    fmt.Println(t) // 2025-11-12 22:14:01 +0700 +07

    // rule.With creates a Timex whose boundaries are evaluated in Vietnam time.
    base := time.Date(2025, 3, 31, 23, 0, 0, 0, time.UTC)
    local := timefy.AdjustTimezone(base, string(timefy.DefaultTimezoneVietnam))
    txVN := rule.With(local)
    fmt.Println(txVN.BeginningOfMonth()) // 2025-04-01 00:00:00 +0700 +07
    // ↑ Note: 23:00 UTC is 06:00 next-day in Vietnam.
    //   So "beginning of month" evaluates to April 1, not March 1.
}
```

---

## 12. Weekday Navigation

The weekday methods return the most recent occurrence of the named day
(looking backward), except for `Sunday()` which returns the upcoming Sunday
(looking forward):

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    // Reference: Wednesday, January 15, 2025
    t := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
    tx := timefy.New(t)

    fmt.Println(tx.Monday())    // 2025-01-13 (most recent Monday)
    fmt.Println(tx.Tuesday())   // 2025-01-14
    fmt.Println(tx.Wednesday()) // 2025-01-15 (today)
    fmt.Println(tx.Thursday())  // 2025-01-09 (most recent Thursday before Wed)
    fmt.Println(tx.Friday())    // 2025-01-10
    fmt.Println(tx.Saturday())  // 2025-01-11
    fmt.Println(tx.Sunday())    // 2025-01-19 (upcoming Sunday)
    fmt.Println(tx.EndOfSunday()) // 2025-01-19 23:59:59.999999999

    // Pass an optional time string to evaluate from a different base:
    fmt.Println(tx.Monday("15:35:34"))  // 2025-01-13 15:35:34
    // The time component from the string is applied to the computed day.

    // Global variants (use time.Now() as base):
    fmt.Println(timefy.Monday())
    fmt.Println(timefy.Sunday())
    fmt.Println(timefy.EndOfSunday())
}
```

---

## 13. Human-Readable Relative Time

### 13.1 TimeAgo

`TimeAgo` returns a natural-language string describing how long ago a past
time occurred.  It uses `time.Now()` as the reference point.

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    cases := []time.Duration{
        30 * time.Second,
        90 * time.Second,
        5 * time.Minute,
        90 * time.Minute,
        3 * time.Hour,
        36 * time.Hour,
        5 * 24 * time.Hour,
        10 * 24 * time.Hour,
        25 * 24 * time.Hour,
        50 * 24 * time.Hour,
        200 * 24 * time.Hour,
        500 * 24 * time.Hour,
        800 * 24 * time.Hour,
    }

    for _, d := range cases {
        past := time.Now().Add(-d)
        fmt.Printf("  %-12v ago → %q\n", d, timefy.New(past).TimeAgo())
    }
    // Output (approximately):
    //   30s         ago → "just now"
    //   1m30s       ago → "1 minute ago"
    //   5m0s        ago → "5 minutes ago"
    //   1h30m0s     ago → "1 hour ago"
    //   3h0m0s      ago → "3 hours ago"
    //   36h0m0s     ago → "1 day ago"
    //   120h0m0s    ago → "5 days ago"
    //   240h0m0s    ago → "1 week ago"
    //   600h0m0s    ago → "3 weeks ago"
    //   1200h0m0s   ago → "1 month ago"
    //   4800h0m0s   ago → "6 months ago"
    //   12000h0m0s  ago → "1 year ago"
    //   19200h0m0s  ago → "2 years ago"
}
```

**Threshold table:**

| Condition | Output |
|-----------|--------|
| `seconds < 60` | `"just now"` |
| `minutes < 2` | `"1 minute ago"` |
| `minutes < 60` | `"N minutes ago"` |
| `hours < 2` | `"1 hour ago"` |
| `hours < 24` | `"N hours ago"` |
| `days < 2` | `"1 day ago"` |
| `days < 7` | `"N days ago"` |
| `days < 14` | `"1 week ago"` |
| `days < 30` | `"N weeks ago"` |
| `days < 60` | `"1 month ago"` |
| `days < 365` | `"N months ago"` |
| `days < 730` | `"1 year ago"` |
| else | `"N years ago"` |

### 13.2 TimeUntil

`TimeUntil` returns a natural-language string describing how long until a
future time arrives.

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    cases := []time.Duration{
        30 * time.Second,
        90 * time.Second,
        10 * time.Minute,
        2 * time.Hour,
        36 * time.Hour,
        5 * 24 * time.Hour,
        15 * 24 * time.Hour,
        45 * 24 * time.Hour,
        200 * 24 * time.Hour,
        800 * 24 * time.Hour,
    }

    for _, d := range cases {
        future := time.Now().Add(d)
        fmt.Printf("  in %-12v → %q\n", d, timefy.New(future).TimeUntil())
    }

    // Past time → "in the past"
    past := time.Now().Add(-1 * time.Hour)
    fmt.Println(timefy.New(past).TimeUntil()) // "in the past"
}
```

---

## 14. Duration Calculations

Compute the integer difference between two times in days, weeks, months,
or years.

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
    tx := timefy.New(start)

    end1 := start.AddDate(0, 0, 45)
    end2 := start.AddDate(0, 0, 21)
    end3 := start.AddDate(0, 3, 0)
    end4 := start.AddDate(2, 6, 0)

    // DurationInDays: absolute hours / 24 (truncated toward zero).
    fmt.Println(tx.DurationInDays(end1))   // -45 (start is 45 days before end1)
    fmt.Println(tx.DurationInDays(start))  // 0

    // DurationInWeeks: DurationInDays / 7.
    fmt.Println(tx.DurationInWeeks(end2))  // -3

    // DurationInMonths: year-and-month difference.
    // Formula: (tx.Year - other.Year) * 12 + (tx.Month - other.Month)
    // A negative value means tx is before other.
    fmt.Println(tx.DurationInMonths(end3)) // -3

    // DurationInYears: tx.Year - other.Year.
    fmt.Println(tx.DurationInYears(end4))  // -2

    // Practical: age calculation
    born := time.Date(1990, 6, 15, 0, 0, 0, 0, time.UTC)
    now := time.Date(2025, 3, 22, 0, 0, 0, 0, time.UTC)
    age := timefy.New(now).DurationInYears(born)
    fmt.Printf("Age: %d years\n", age) // 34
}
```

**Important caveat:**
`DurationInDays` uses `t.Sub(other).Hours() / 24`, which measures wall-clock
hours.  This can differ from calendar days across DST boundaries.  For
calendar-day semantics across DST, use `AddDate(0, 0, n)` and count manually.

---

## 15. Time Comparisons

### 15.1 Ordinal Comparisons

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    t := time.Date(2025, 3, 15, 12, 0, 0, 0, time.UTC)
    tx := timefy.New(t)

    earlier := time.Date(2025, 3, 10, 12, 0, 0, 0, time.UTC)
    later   := time.Date(2025, 3, 20, 12, 0, 0, 0, time.UTC)

    fmt.Println(tx.IsBefore(later))    // true
    fmt.Println(tx.IsBefore(earlier))  // false
    fmt.Println(tx.IsAfter(earlier))   // true
    fmt.Println(tx.IsAfter(later))     // false
}
```

### 15.2 Range Membership

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    t := time.Date(2025, 3, 15, 12, 0, 0, 0, time.UTC)
    tx := timefy.New(t)

    start := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
    end   := time.Date(2025, 3, 31, 23, 59, 59, 999999999, time.UTC)

    // IsBetween is inclusive on both endpoints.
    // t == start → true; t == end → true.
    fmt.Println(tx.IsBetween(start, end)) // true

    outsideEnd := time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)
    fmt.Println(tx.IsBetween(end, outsideEnd)) // false (t is before start)
}
```

### 15.3 Calendrical Identity

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    t  := time.Date(2025, 3, 15, 8, 0, 0, 0, time.UTC)
    t2 := time.Date(2025, 3, 15, 22, 0, 0, 0, time.UTC)
    t3 := time.Date(2025, 3, 16, 8, 0, 0, 0, time.UTC)
    t4 := time.Date(2025, 4, 15, 8, 0, 0, 0, time.UTC)
    t5 := time.Date(2026, 3, 15, 8, 0, 0, 0, time.UTC)

    tx := timefy.New(t)

    fmt.Println(tx.IsSameDay(t2))   // true  (same year-day)
    fmt.Println(tx.IsSameDay(t3))   // false
    fmt.Println(tx.IsSameMonth(t4)) // false (different month)
    fmt.Println(tx.IsSameMonth(t3)) // true  (same year-month)
    fmt.Println(tx.IsSameYear(t5))  // false
    fmt.Println(tx.IsSameYear(t4))  // true
}
```

### 15.4 Temporal Position

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    now := time.Now()

    fmt.Println(timefy.New(now).IsToday())                      // true
    fmt.Println(timefy.New(now.AddDate(0, 0, -1)).IsYesterday()) // true
    fmt.Println(timefy.New(now.AddDate(0, 0, 1)).IsTomorrow())   // true
    fmt.Println(timefy.New(now.Add(-1 * time.Hour)).IsPast())    // true
    fmt.Println(timefy.New(now.Add(1 * time.Hour)).IsFuture())   // true
}
```

### 15.5 Day-Type Checks

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    saturday := time.Date(2025, 1, 18, 0, 0, 0, 0, time.UTC)
    sunday   := time.Date(2025, 1, 19, 0, 0, 0, 0, time.UTC)
    monday   := time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC)
    friday   := time.Date(2025, 1, 17, 0, 0, 0, 0, time.UTC)

    fmt.Println(timefy.New(saturday).IsWeekend()) // true
    fmt.Println(timefy.New(sunday).IsWeekend())   // true
    fmt.Println(timefy.New(monday).IsWeekend())   // false
    fmt.Println(timefy.New(friday).IsWeekend())   // false

    fmt.Println(timefy.New(saturday).IsWeekday()) // false
    fmt.Println(timefy.New(monday).IsWeekday())   // true
}
```

---

## 16. Between (String-based Range Check)

`Between` parses two time strings and checks whether the current `Timex`
time falls between them:

```go
package main

import (
    "fmt"

    "github.com/sivaosorg/timefy"
)

func main() {
    // Global function: checks if time.Now() is between the two strings.
    fmt.Println(timefy.Between("2020-01-01", "2030-12-31")) // true
    fmt.Println(timefy.Between("2030-01-01", "2031-12-31")) // false

    // Timex method: checks if tx.Time is between the strings.
    // Useful for checking a known time against a window.
    // (Uses MustParse internally — panics if strings are unparseable.)
}
```

---

## 17. Elapsed Time Measurements

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    // Measure how long since a known point.
    start := time.Date(2025, 3, 22, 8, 0, 0, 0, time.UTC)

    // These all use time.Since(start) internally.
    fmt.Printf("Hours since start:   %.2f\n", timefy.SinceHour(start))
    fmt.Printf("Minutes since start: %.2f\n", timefy.SinceMinute(start))
    fmt.Printf("Seconds since start: %.2f\n", timefy.SinceSecond(start))

    // Typical use: performance monitoring.
    taskStart := time.Now()
    // ... perform task ...
    elapsed := timefy.SinceSecond(taskStart)
    if elapsed > 5.0 {
        fmt.Printf("Warning: task took %.2fs (> 5s SLA)\n", elapsed)
    }
}
```

---

## 18. Leap Year Utilities

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    // IsLeapYear(year): checks the standard Gregorian rules:
    //   - divisible by 4   AND
    //   - NOT divisible by 100, UNLESS also divisible by 400.
    fmt.Println(timefy.IsLeapYear(2024)) // true  (÷4, not ÷100)
    fmt.Println(timefy.IsLeapYear(2023)) // false
    fmt.Println(timefy.IsLeapYear(1900)) // false (÷100 but not ÷400)
    fmt.Println(timefy.IsLeapYear(2000)) // true  (÷400)
    fmt.Println(timefy.IsLeapYear(2100)) // false (÷100 but not ÷400)

    // IsLeapYearN(t): convenience wrapper extracting the year from t.
    t2024 := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
    fmt.Println(timefy.IsLeapYearN(t2024)) // true

    // Practical: determine February end for a dynamic billing calculation.
    for _, yr := range []int{2023, 2024, 2025} {
        feb := time.Date(yr, time.February, 1, 0, 0, 0, 0, time.UTC)
        endFeb := timefy.New(feb).EndOfMonth()
        fmt.Printf("  %d Feb ends on day %d (leap=%v)\n",
            yr, endFeb.Day(), timefy.IsLeapYear(yr))
    }
    // Output:
    //   2023 Feb ends on day 28 (leap=false)
    //   2024 Feb ends on day 29 (leap=true)
    //   2025 Feb ends on day 28 (leap=false)
}
```

---

## 19. Weekdays in Range

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    // Returns all Monday–Friday dates in [start, end] (inclusive).
    start := time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC) // Monday
    end   := time.Date(2025, 1, 19, 0, 0, 0, 0, time.UTC) // Sunday

    weekdays := timefy.GetWeekdaysInRange(start, end)
    for _, d := range weekdays {
        fmt.Println(d.Format("Mon 2006-01-02"))
    }
    // Mon 2025-01-13
    // Tue 2025-01-14
    // Wed 2025-01-15
    // Thu 2025-01-16
    // Fri 2025-01-17

    fmt.Printf("Count: %d\n", len(weekdays)) // 5
}
```

**Known limitation:** `GetWeekdaysInRange` has a conditional branch that, when
the date range falls within a leap year, only counts February 29 as a valid
weekday and skips all other weekdays in that year.  For date ranges that fall
entirely within a non-leap year, the function works correctly.  When working
with date ranges that span any part of a leap year, use `IsWeekday()` in a
manual iteration loop instead:

```go
// Portable alternative for any year:
var weekdays []time.Time
for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
    if timefy.New(d).IsWeekday() {
        weekdays = append(weekdays, d)
    }
}
```

---

## 20. Tolerance Window

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    now := time.Now()

    // IsWithinTolerance: returns true if t is within ±1 minute of now.
    // Useful for idempotency checks (e.g., webhook replays).
    fmt.Println(timefy.IsWithinTolerance(now))                          // true
    fmt.Println(timefy.IsWithinTolerance(now.Add(30 * time.Second)))    // true
    fmt.Println(timefy.IsWithinTolerance(now.Add(-45 * time.Second)))   // true
    fmt.Println(timefy.IsWithinTolerance(now.Add(2 * time.Minute)))     // false
    fmt.Println(timefy.IsWithinTolerance(now.Add(-2 * time.Minute)))    // false
}
```

---

## 21. Quarter Number

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    // Package-level: quarter of time.Now().
    fmt.Println(timefy.Quarter()) // e.g., 1 for January–March

    // Timex method: quarter of any time.Time.
    months := []time.Month{
        time.January, time.March,    // Q1
        time.April,   time.June,     // Q2
        time.July,    time.September,// Q3
        time.October, time.December, // Q4
    }
    for _, m := range months {
        t := time.Date(2025, m, 1, 0, 0, 0, 0, time.UTC)
        fmt.Printf("  %s → Q%d\n", m, timefy.New(t).Quarter())
    }
}
```

---

## 22. Rule Configuration — Full API

`Rule` is the configuration builder for `Timex`.  All methods return
`*Rule` for fluent chaining.

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

func main() {
    // NewRule() creates a Rule with defaults:
    //   weekStartDay = time.Sunday
    //   timeFormats  = TimeFormats (40+ built-in patterns)
    //   timeLocation = nil (uses the time value's location)
    rule := timefy.NewRule()

    // WithWeekStartDay: change which day starts the week.
    rule.WithWeekStartDay(time.Monday)

    // WithTimeFormats: REPLACE the format registry entirely.
    rule.WithTimeFormats([]string{
        "2006-01-02 15:04:05",
        "2006-01-02",
    })

    // AppendTimeFormat: ADD custom formats to the existing registry.
    rule.AppendTimeFormat("02 Jan 2006 15:04", "Jan 02, 2006")

    // AppendTimeFormatRFC: add a typed TimeFormatRFC constant.
    rule.AppendTimeFormatRFC(
        timefy.TimeFormat20060102T150405,
        timefy.TimeFormat20060102T150405999999,
    )

    // AppendTimeFormatRFCshort: add a typed TimeRFC constant.
    rule.AppendTimeFormatRFCshort(timefy.TimeRFC01T150405)

    // WithLocation: set timezone using a *time.Location.
    rule.WithLocation(time.UTC)

    // WithLocationRFC: set timezone using a ZoneRFC constant.
    rule.WithLocationRFC(timefy.DefaultTimezoneVietnam)

    // ApplyUTCLoc: shorthand for WithLocation(time.UTC).
    rule.ApplyUTCLoc()

    // ApplyLocalLoc: shorthand for WithLocation(time.Local).
    rule.ApplyLocalLoc()

    // After configuring the rule, use it to:

    // 1. Wrap a time.Time into a Timex:
    tx := rule.With(time.Now())
    fmt.Println(tx.BeginningOfWeek()) // uses Monday as week start

    // 2. Parse a string using the rule's formats and location:
    t, err := rule.Parse("2025-03-15 14:30:00")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(t)

    // 3. MustParse — panics on failure:
    t2 := rule.MustParse("2025-03-15")
    fmt.Println(t2)

    // Chaining example (all on one line):
    specialRule := timefy.NewRule().
        WithWeekStartDay(time.Monday).
        WithLocationRFC(timefy.DefaultTimezoneLondon).
        AppendTimeFormat("02/01/2006")

    txSpecial := specialRule.With(time.Date(2025, 3, 19, 0, 0, 0, 0, time.UTC))
    fmt.Println(txSpecial.BeginningOfWeek()) // Monday 2025-03-17
}
```

---

## 23. Advanced Patterns

### 23.1 Fluent Builder Chains

Build reusable rule instances at package initialization:

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

var (
    // ruleISO: ISO 8601 week (Monday-start), UTC.
    ruleISO = timefy.NewRule().
        WithWeekStartDay(time.Monday).
        ApplyUTCLoc()

    // ruleVN: Vietnam timezone, Monday-start.
    ruleVN = timefy.NewRule().
        WithWeekStartDay(time.Monday).
        WithLocationRFC(timefy.DefaultTimezoneVietnam)

    // ruleStrict: only accept the exact format "2006-01-02 15:04:05".
    ruleStrict = timefy.NewRule().
        WithTimeFormats([]string{"2006-01-02 15:04:05"}).
        ApplyUTCLoc()
)

func main() {
    t := time.Date(2025, 3, 19, 0, 0, 0, 0, time.UTC)

    fmt.Println(ruleISO.With(t).BeginningOfWeek())   // 2025-03-17 (Monday)
    fmt.Println(ruleVN.With(t).BeginningOfMonth())   // 2025-03-01 in +07

    v, err := ruleStrict.Parse("2025-03-15 14:30:00")
    fmt.Println(v, err) // ok

    _, err = ruleStrict.Parse("2025-03-15") // fails: no matching format
    fmt.Println(err) // can't parse string as time: 2025-03-15
}
```

### 23.2 Multi-tenant Timezone Isolation

In SaaS applications each tenant may operate in a different timezone.
Carry the rule as part of a tenant context:

```go
package main

import (
    "fmt"
    "time"

    "github.com/sivaosorg/timefy"
)

type TenantCtx struct {
    ID       string
    Timezone timefy.ZoneRFC
    Rule     *timefy.Rule
}

func newTenantCtx(id string, tz timefy.ZoneRFC) TenantCtx {
    return TenantCtx{
        ID:       id,
        Timezone: tz,
        Rule:     timefy.NewRule().WithLocationRFC(tz),
    }
}

func billingWindowFor(tc TenantCtx, t time.Time) (start, end time.Time) {
    // Convert the UTC time to the tenant's local timezone first,
    // then apply the rule-bound Timex methods so that month boundaries
    // are evaluated in the tenant's local calendar context.
    local := timefy.AdjustTimezone(t, string(tc.Timezone))
    tx := tc.Rule.With(local)
    return tx.BeginningOfMonth(), tx.EndOfMonth()
}

func main() {
    tenants := []TenantCtx{
        newTenantCtx("tenant-vn", timefy.DefaultTimezoneVietnam),
        newTenantCtx("tenant-us", timefy.DefaultTimezoneNewYork),
        newTenantCtx("tenant-uk", timefy.DefaultTimezoneLondon),
    }

    // A UTC time that straddles midnight in different timezones.
    utcBorderTime := time.Date(2025, 3, 31, 23, 30, 0, 0, time.UTC)

    for _, tc := range tenants {
        s, e := billingWindowFor(tc, utcBorderTime)
        fmt.Printf("  [%s] billing window: %s → %s\n",
            tc.ID, s.Format("Jan 02"), e.Format("Jan 02"))
    }
    // [tenant-vn] billing window: Apr 01 → Apr 30  (23:30 UTC = 06:30 next-day in +07)
    // [tenant-us] billing window: Mar 01 → Mar 31  (23:30 UTC = 19:30 same-day in EDT)
    // [tenant-uk] billing window: Mar 01 → Mar 31  (23:30 UTC = 23:30 same-day in GMT)
}
```

### 23.3 Custom Format Registries

Override the global `TimeFormats` in `main()` or use per-rule formats:

```go
package main

import (
    "fmt"

    "github.com/sivaosorg/timefy"
)

func init() {
    // Register custom formats globally (affects Parse, MustParse, etc.).
    timefy.TimeFormats = append(timefy.TimeFormats,
        "02 Jan 2006 15:04",     // British short form
        "January 2, 2006",       // Long English form
        "2 January 2006",        // Reversed long form
    )
}

func main() {
    t, _ := timefy.Parse("15 Mar 2025 14:30")
    fmt.Println(t) // 2025-03-15 14:30:00

    t2, _ := timefy.Parse("March 15, 2025")
    fmt.Println(t2) // 2025-03-15 00:00:00
}
```

### 23.4 Safe vs Panic Patterns

| Scenario | Function | Behavior |
|----------|----------|----------|
| User input (HTTP, CLI) | `Parse`, `ParseInLocation` | Returns `error` — handle it |
| Config file constants | `MustParse`, `MustParseInLocation` | Panics — appropriate for startup |
| Trusted internal code | `MustParse` | Panics on programmer error |
| Dynamic timezone from DB | `SetTimezone` | Returns `error` if TZ invalid |
| Trusted timezone constant | `AdjustTimezone` | Silent fallback on error |

```go
// ✅ Safe pattern for HTTP handlers:
func handleRequest(dateStr string) (time.Time, error) {
    t, err := timefy.Parse(dateStr)
    if err != nil {
        return time.Time{}, fmt.Errorf("invalid date %q: %w", dateStr, err)
    }
    return t, nil
}

// ✅ Panic pattern for program constants:
var startOfEpoch = timefy.MustParse("2020-01-01")
```

---

## 24. Enterprise Use Cases

The `examples/` directory contains five complete, runnable enterprise programs
that demonstrate production-grade patterns.

### `examples/billing/` — SaaS Subscription Billing

**Business problem:** A multi-tenant SaaS platform bills customers monthly,
quarterly, and annually.  Customers sign up mid-period and receive pro-rated
charges.  Some plans include a free trial period.

**Key timefy APIs used:**

| API | Purpose |
|-----|---------|
| `New(t).BeginningOfMonth()` | Start of billing period |
| `New(t).EndOfMonth()` | End of billing period |
| `New(t).BeginningOfQuarter()` | Q1/Q2/Q3/Q4 start |
| `New(t).EndOfQuarter()` | Q1/Q2/Q3/Q4 end |
| `New(t).BeginningOfYear()` | Annual period start |
| `New(t).EndOfYear()` | Annual period end |
| `AddDay(t, n)` | Trial expiry calculation |
| `AdjustTimezone` | Tenant-local boundary evaluation |
| `IsLeapYear` | February billing edge case |
| `IsBetween` | Trial membership test |
| `DurationInDays` | Pro-rate day counting |

```bash
cd examples/billing && go run main.go
```

### `examples/global_team/` — Multi-Timezone Team Coordination

**Business problem:** An engineering organization spans five cities.  The
platform must find meeting overlap windows, detect business-hour availability,
and generate cross-team working-day schedules that respect public holidays.

**Key timefy APIs used:**

| API | Purpose |
|-----|---------|
| `AdjustTimezone` | Convert UTC probe times to local |
| `New(t).IsWeekday()` | Exclude weekends |
| `AddHour(t, 1)` | Scan 1-hour overlap windows |
| `AddDay(t, 1)` | Advance day cursor |
| `New(t).IsSameDay` | Holiday matching |
| `New(t).BeginningOfQuarter()` | Quarter display |
| `New(t).TimeAgo()` / `TimeUntil()` | Release timeline display |

```bash
cd examples/global_team && go run main.go
```

### `examples/audit/` — Enterprise Audit Trail & Compliance

**Business problem:** A financial services platform records all data-access
events, enforces GDPR Article 33 (72-hour breach notification), detects SLA
breaches for support tickets, and generates CISO reports.

**Key timefy APIs used:**

| API | Purpose |
|-----|---------|
| `AddHour(t, 72)` | GDPR 72-hour notification deadline |
| `AddHour(t, 1)` | Business-hour SLA clock ticking |
| `New(t).IsWeekday()` | SLA clock excludes weekends |
| `New(t).BeginningOfHour()` | Hourly event bucketing |
| `New(t).IsBetween` | Events-in-window query |
| `New(t).BeginningOfWeek()` / `EndOfWeek()` | Weekly report window |
| `New(t).BeginningOfQuarter()` / `Quarter()` | Quarterly KPI window |
| `BeginOfDay` / `AddDay` | 90-day retention threshold |
| `FormatTimex` | Structured log timestamp extraction |
| `New(t).TimeAgo()` | Human-readable event age |

```bash
cd examples/audit && go run main.go
```

### `examples/reporting/` — Financial Quarter & Fiscal-Year Reporting

**Business problem:** A publicly-traded company needs precise SEC filing
windows, YTD periods, QoQ comparisons, and a July-1-start fiscal year
calendar.

**Key timefy APIs used:**

| API | Purpose |
|-----|---------|
| `New(t).BeginningOfQuarter()` | Calendar Q1–Q4 start |
| `New(t).EndOfQuarter()` | Calendar Q1–Q4 end |
| `New(t).BeginningOfYear()` | YTD start |
| `New(t).EndOfYear()` | Year-end boundary |
| `New(t).BeginningOfHalf()` | H1/H2 start |
| `AddDay(t, 45)` | 10-Q filing deadline |
| `AddDay(t, 60)` | 10-K filing deadline |
| `New(t).DurationInDays` | Days-until-deadline countdown |
| `GetWeekdaysInRange` | Trading-day count |
| `IsLeapYear` | February calendar check |
| `New(t).Quarter()` | Quarter number label |

```bash
cd examples/reporting && go run main.go
```

### `examples/scheduler/` — Business-Day-Aware Job Scheduler

**Business problem:** A financial data platform schedules ETL jobs on daily,
monthly, and quarterly boundaries, ensuring jobs fire only on business days and
are not duplicated within a configurable time window.

**Key timefy APIs used:**

| API | Purpose |
|-----|---------|
| `New(t).BeginningOfMonth()` | Monthly job anchor |
| `New(t).BeginningOfQuarter()` | Quarterly job anchor |
| `New(t).IsWeekday()` | Skip weekends |
| `AddDay(t, 1)` | Forward-step day cursor |
| `AddHour(t, 1)` | Business-hour SLA tick |
| `AddMinute(t, -n)` | Anti-duplication window start |
| `New(t).IsBetween` | Duplicate-run detection |
| `GetWeekdaysInRange` | Week overview rendering |
| `New(t).Quarter()` | Quarter label on schedule output |
| `New(t).BeginningOfWeek()` / `EndOfWeek()` | Week schedule bounds |
| `New(t).TimeAgo()` / `TimeUntil()` | Last-ran / next-run display |

```bash
cd examples/scheduler && go run main.go
```

---

## 25. API Quick-Reference

### Package-level functions

| Function | Returns | Description |
|----------|---------|-------------|
| `With(t)` | `*Timex` | Wrap time.Time in Timex |
| `New(t)` | `*Timex` | Alias for With |
| `BeginningOfMinute()` | `time.Time` | Start of current minute |
| `BeginningOfHour()` | `time.Time` | Start of current hour |
| `BeginningOfDay()` | `time.Time` | Midnight today |
| `BeginningOfWeek()` | `time.Time` | Start of current week |
| `BeginningOfMonth()` | `time.Time` | 1st of current month |
| `BeginningOfQuarter()` | `time.Time` | 1st of current quarter |
| `BeginningOfYear()` | `time.Time` | Jan 1 of current year |
| `EndOfMinute()` | `time.Time` | 59.999999999s of current minute |
| `EndOfHour()` | `time.Time` | 59:59.999999999 of current hour |
| `EndOfDay()` | `time.Time` | 23:59:59.999999999 today |
| `EndOfWeek()` | `time.Time` | Last ns of current week |
| `EndOfMonth()` | `time.Time` | Last ns of current month |
| `EndOfQuarter()` | `time.Time` | Last ns of current quarter |
| `EndOfYear()` | `time.Time` | Dec 31 23:59:59.999999999 |
| `BeginOfDay(t)` | `time.Time` | Midnight of date t |
| `FEndOfDay(t)` | `time.Time` | 23:59:59 of date t |
| `PrevBeginOfDay(t, n)` | `time.Time` | Midnight n days before t |
| `PrevEndOfDay(t, n)` | `time.Time` | 23:59:59 n days before t |
| `AddSecond(t, n)` | `time.Time` | t ± n seconds |
| `AddMinute(t, n)` | `time.Time` | t ± n minutes |
| `AddHour(t, n)` | `time.Time` | t ± n hours |
| `AddDay(t, n)` | `time.Time` | t ± n×24h |
| `Monday(s...)` | `time.Time` | Most recent Monday |
| `Tuesday(s...)` | `time.Time` | Most recent Tuesday |
| `Wednesday(s...)` | `time.Time` | Most recent Wednesday |
| `Thursday(s...)` | `time.Time` | Most recent Thursday |
| `Friday(s...)` | `time.Time` | Most recent Friday |
| `Saturday(s...)` | `time.Time` | Most recent Saturday |
| `Sunday(s...)` | `time.Time` | Upcoming Sunday |
| `EndOfSunday()` | `time.Time` | End of upcoming Sunday |
| `Quarter()` | `uint` | Quarter of today (1–4) |
| `Parse(s...)` | `(time.Time, error)` | Parse string(s) |
| `MustParse(s...)` | `time.Time` | Parse or panic |
| `ParseInLocation(loc, s...)` | `(time.Time, error)` | Parse in timezone |
| `MustParseInLocation(loc, s...)` | `time.Time` | Parse in timezone or panic |
| `Between(t1, t2)` | `bool` | Is now between t1 and t2? |
| `FormatRFC(layout)` | `string` | Format now with layout |
| `FormatRFCshort(layout)` | `string` | Format now with short layout |
| `FFormatRFC(t, layout)` | `string` | Format t with layout |
| `FFormatRFCshort(t, layout)` | `string` | Format t with short layout |
| `DefaultFormatRFC()` | `string` | Format now (microseconds) |
| `FDefaultFormatRFC(t)` | `string` | Format t (microseconds) |
| `SetTimezone(t, tz)` | `(time.Time, error)` | Convert to timezone |
| `AdjustTimezone(t, tz)` | `time.Time` | Convert; fallback on error |
| `IsLeapYear(year)` | `bool` | Is year a leap year? |
| `IsLeapYearN(t)` | `bool` | Is t.Year() a leap year? |
| `IsWithinTolerance(t)` | `bool` | Is t within ±1 minute of now? |
| `SinceHour(t)` | `float64` | Hours since t |
| `SinceMinute(t)` | `float64` | Minutes since t |
| `SinceSecond(t)` | `float64` | Seconds since t |
| `GetWeekdaysInRange(s, e)` | `[]time.Time` | Mon–Fri dates in [s,e] |
| `FormatTimex(t)` | `[]int` | [ns, sec, min, hr, day, mo, yr] |
| `NewRule()` | `*Rule` | New config with defaults |

### Timex methods

| Method | Returns | Description |
|--------|---------|-------------|
| `.BeginningOfMinute()` | `time.Time` | Start of minute |
| `.BeginningOfHour()` | `time.Time` | Start of hour |
| `.BeginningOfDay()` | `time.Time` | Midnight |
| `.BeginningOfWeek()` | `time.Time` | Start of week (per weekStartDay) |
| `.BeginningOfMonth()` | `time.Time` | 1st of month |
| `.BeginningOfQuarter()` | `time.Time` | 1st of quarter |
| `.BeginningOfHalf()` | `time.Time` | 1st of half-year |
| `.BeginningOfYear()` | `time.Time` | Jan 1st |
| `.EndOfMinute()` | `time.Time` | Last ns of minute |
| `.EndOfHour()` | `time.Time` | Last ns of hour |
| `.EndOfDay()` | `time.Time` | Last ns of day |
| `.EndOfWeek()` | `time.Time` | Last ns of week |
| `.EndOfMonth()` | `time.Time` | Last ns of month |
| `.EndOfQuarter()` | `time.Time` | Last ns of quarter |
| `.EndOfHalf()` | `time.Time` | Last ns of half-year |
| `.EndOfYear()` | `time.Time` | Last ns of year |
| `.Monday(s...)` | `time.Time` | Most recent Monday |
| `.Tuesday(s...)` | `time.Time` | Most recent Tuesday |
| `.Wednesday(s...)` | `time.Time` | Most recent Wednesday |
| `.Thursday(s...)` | `time.Time` | Most recent Thursday |
| `.Friday(s...)` | `time.Time` | Most recent Friday |
| `.Saturday(s...)` | `time.Time` | Most recent Saturday |
| `.Sunday(s...)` | `time.Time` | Upcoming Sunday |
| `.EndOfSunday()` | `time.Time` | End of upcoming Sunday |
| `.Quarter()` | `uint` | Quarter (1–4) |
| `.Parse(s...)` | `(time.Time, error)` | Parse with rule context |
| `.MustParse(s...)` | `time.Time` | Parse or panic |
| `.Between(begin, end)` | `bool` | Is tx.Time between begin/end? |
| `.FormatRFC(layout)` | `string` | Format with full layout |
| `.FormatRFCshort(layout)` | `string` | Format with short layout |
| `.DefaultFormatRFC()` | `string` | Format (microseconds) |
| `.TimeAgo()` | `string` | Human-readable past time |
| `.TimeUntil()` | `string` | Human-readable future time |
| `.DurationInDays(other)` | `int` | Days between tx and other |
| `.DurationInWeeks(other)` | `int` | Weeks between tx and other |
| `.DurationInMonths(other)` | `int` | Months between tx and other |
| `.DurationInYears(other)` | `int` | Years between tx and other |
| `.IsWeekend()` | `bool` | Sat or Sun? |
| `.IsWeekday()` | `bool` | Mon–Fri? |
| `.IsBefore(other)` | `bool` | tx < other? |
| `.IsAfter(other)` | `bool` | tx > other? |
| `.IsBetween(start, end)` | `bool` | start ≤ tx ≤ end? |
| `.IsSameDay(other)` | `bool` | Same calendar day? |
| `.IsSameMonth(other)` | `bool` | Same year-month? |
| `.IsSameYear(other)` | `bool` | Same year? |
| `.IsToday()` | `bool` | tx.Date == today? |
| `.IsYesterday()` | `bool` | tx.Date == yesterday? |
| `.IsTomorrow()` | `bool` | tx.Date == tomorrow? |
| `.IsPast()` | `bool` | tx < now? |
| `.IsFuture()` | `bool` | tx > now? |

### Rule methods

| Method | Returns | Description |
|--------|---------|-------------|
| `NewRule()` | `*Rule` | Default config |
| `.WithWeekStartDay(d)` | `*Rule` | Set week-start day |
| `.WithTimeFormats(fs)` | `*Rule` | Replace format list |
| `.WithLocation(loc)` | `*Rule` | Set *time.Location |
| `.WithLocationRFC(z)` | `*Rule` | Set timezone by ZoneRFC |
| `.ApplyUTCLoc()` | `*Rule` | Set UTC |
| `.ApplyLocalLoc()` | `*Rule` | Set local |
| `.AppendTimeFormat(s...)` | `*Rule` | Append string formats |
| `.AppendTimeFormatRFC(f...)` | `*Rule` | Append TimeFormatRFC |
| `.AppendTimeFormatRFCshort(f...)` | `*Rule` | Append TimeRFC |
| `.With(t)` | `*Timex` | Create Timex with this rule |
| `.Parse(s...)` | `(time.Time, error)` | Parse with this rule |
| `.MustParse(s...)` | `time.Time` | Parse or panic |

---

## 26. Format Constant Reference

### Short formats (TimeRFC)

| Constant | Pattern | Example |
|----------|---------|---------|
| `TimeRFC01T150405` | `15:04:05` | `13:45:30` |
| `TimeRFC02D150405` | `15.04.05` | `13.45.30` |
| `TimeRFC03H150405` | `15-04-05` | `13-45-30` |
| `TimeRFC04M150405` | `15/04/05` | `13/45/30` |
| `TimeRFC05S150405` | `15 04 05` | `13 45 30` |
| `TimeRFC06T150405` | `150405` | `134530` |
| `TimeRFC07D150405` | `15:04` | `13:45` |
| `TimeRFC08H150405` | `15.04` | `13.45` |
| `TimeRFC09M150405` | `15-04` | `13-45` |
| `TimeRFC10S150405` | `15/04` | `13/45` |
| `TimeRFC11T150405` | `1504` | `1345` |

### Full formats (TimeFormatRFC) — selected

| Constant | Pattern | Example |
|----------|---------|---------|
| `TimeFormat20060102T150405999999` | `2006-01-02T15:04:05.999999` | `2023-08-15T13:45:30.123456` |
| `TimeFormat20060102T150405` | `2006-01-02T15:04:05` | `2023-08-15T13:45:30` |
| `TimeFormat20060102150405` | `2006-01-02 15:04:05` | `2023-08-15 13:45:30` |
| `TimeFormat20060102150405999999` | `2006-01-02 15:04:05.999999` | `2023-08-15 13:45:30.123456` |
| `TimeFormat20060102` | `2006-01-02` | `2023-08-15` |
| `TimeFormat200601` | `2006-01` | `2023-08` |
| `TimeFormat20060102150405Z0700` | `2006-01-02 15:04:05 -07:00` | `2023-08-15 13:45:30 -05:00` |
| `TimeFormat20060102150405Z0700RFC3339` | `2006-01-02T15:04:05-07:00` | `2023-08-15T13:45:30-05:00` |

---

## 27. Timezone Constant Reference

| Constant | IANA Name | UTC Offset (standard) |
|----------|-----------|----------------------|
| `DefaultTimezoneVietnam` | `Asia/Ho_Chi_Minh` | +07:00 |
| `DefaultTimezoneNewYork` | `America/New_York` | -05:00 / -04:00 (DST) |
| `DefaultTimezoneLondon` | `Europe/London` | +00:00 / +01:00 (BST) |
| `DefaultTimezoneTokyo` | `Asia/Tokyo` | +09:00 |
| `DefaultTimezoneSydney` | `Australia/Sydney` | +10:00 / +11:00 (AEDT) |
| `DefaultTimezoneParis` | `Europe/Paris` | +01:00 / +02:00 (CEST) |
| `DefaultTimezoneMoscow` | `Europe/Moscow` | +03:00 |
| `DefaultTimezoneLosAngeles` | `America/Los_Angeles` | -08:00 / -07:00 (PDT) |
| `DefaultTimezoneManila` | `Asia/Manila` | +08:00 |
| `DefaultTimezoneKualaLumpur` | `Asia/Kuala_Lumpur` | +08:00 |
| `DefaultTimezoneJakarta` | `Asia/Jakarta` | +07:00 |
| `DefaultTimezoneYangon` | `Asia/Yangon` | +06:30 |
| `DefaultTimezoneAuckland` | `Pacific/Auckland` | +12:00 / +13:00 (NZDT) |
| `DefaultTimezoneBangkok` | `Asia/Bangkok` | +07:00 |
| `DefaultTimezoneDelhi` | `Asia/Kolkata` | +05:30 |
| `DefaultTimezoneDubai` | `Asia/Dubai` | +04:00 |
| `DefaultTimezoneCairo` | `Africa/Cairo` | +02:00 |
| `DefaultTimezoneAthens` | `Europe/Athens` | +02:00 / +03:00 (EEST) |
| `DefaultTimezoneRome` | `Europe/Rome` | +01:00 / +02:00 (CEST) |
| `DefaultTimezoneJohannesburg` | `Africa/Johannesburg` | +02:00 |
| `DefaultTimezoneStockholm` | `Europe/Stockholm` | +01:00 / +02:00 (CEST) |
| `DefaultTimezoneOslo` | `Europe/Oslo` | +01:00 / +02:00 (CEST) |
| `DefaultTimezoneHelsinki` | `Europe/Helsinki` | +02:00 / +03:00 (EEST) |
| `DefaultTimezoneKiev` | `Europe/Kiev` | +02:00 / +03:00 (EEST) |
| `DefaultTimezoneBeijing` | `Asia/Shanghai` | +08:00 |
| `DefaultTimezoneSingapore` | `Asia/Singapore` | +08:00 |
| `DefaultTimezoneIslamabad` | `Asia/Karachi` | +05:00 |
| `DefaultTimezoneColombo` | `Asia/Colombo` | +05:30 |
| `DefaultTimezoneDhaka` | `Asia/Dhaka` | +06:00 |
| `DefaultTimezoneKathmandu` | `Asia/Kathmandu` | +05:45 |
| `DefaultTimezoneBrisbane` | `Australia/Brisbane` | +10:00 |
| `DefaultTimezoneWellington` | `Pacific/Auckland` | +12:00 / +13:00 (NZDT) |
| `DefaultTimezonePortMoresby` | `Pacific/Port_Moresby` | +10:00 |
| `DefaultTimezoneSuva` | `Pacific/Fiji` | +12:00 |

---
