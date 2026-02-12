# timefy

**timefy** is a comprehensive Go library designed to enhance productivity by offering a wide range of time-related utilities. Whether you're dealing with date formatting, time zone conversions, or scheduling. Timefy provides a robust toolkit to simplify time management tasks in Go applications.

### Requirements

- Go version 1.23 or higher

### Installation

To install, you can use the following commands based on your preference:

- For a specific version:

  ```bash
  go get github.com/sivaosorg/timefy@v0.0.1
  ```

- For the latest version:
  ```bash
  go get -u github.com/sivaosorg/timefy@latest
  ```

### Getting started

#### Getting timefy

With [Go's module support](https://go.dev/wiki/Modules#how-to-use-modules), `go [build|run|test]` automatically fetches the necessary dependencies when you add the import in your code:

```go
import "github.com/sivaosorg/timefy"
```

#### Usage

> Calculating the time based on current time

```go
package main

import (
	"fmt"
	"time"

	"github.com/sivaosorg/timefy"
)

func main() {
	// time based on current time
	tx := timefy.With(time.Now())        // 2024-10-25 21:31:00.690342
	fmt.Println(tx.BeginningOfMinute())  // 2024-10-25 21:31:00
	fmt.Println(tx.BeginningOfHour())    // 2024-10-25 21:00:00
	fmt.Println(tx.BeginningOfDay())     // 2024-10-25 00:00:00
	fmt.Println(tx.BeginningOfWeek())    // 2024-10-20 00:00:00
	fmt.Println(tx.BeginningOfMonth())   // 2024-10-01 00:00:00
	fmt.Println(tx.BeginningOfQuarter()) // 2024-10-01 00:00:00
	fmt.Println(tx.BeginningOfYear())    // 2024-01-01 00:00:00
	fmt.Println(tx.EndOfMinute())        // 2024-10-25 21:31:59.999999999 +0700 +07
	fmt.Println(tx.EndOfHour())          // 2024-10-25 21:59:59.999999999 +0700 +07
	fmt.Println(tx.EndOfDay())           // 2024-10-25 23:59:59.999999999 +0700 +07
	fmt.Println(tx.EndOfWeek())          // 2024-10-26 23:59:59.999999999 +0700 +07
	fmt.Println(tx.EndOfMonth())         // 2024-10-31 23:59:59.999999999 +0700 +07
	fmt.Println(tx.EndOfQuarter())       // 2024-12-31 23:59:59.999999999 +0700 +07
	fmt.Println(tx.EndOfYear())          // 2024-12-31 23:59:59.999999999 +0700 +07
	tx.WeekStartDay = time.Tuesday       // Default is Sunday
	fmt.Println(tx.EndOfWeek())          // 2024-10-28 23:59:59.999999999 +0700 +07
}
```

> Calculating the time based on another time

```go
package main

import (
	"fmt"
	"time"

	"github.com/sivaosorg/timefy"
)

func main() {
	// time based on another time
	t := time.Date(2025, 01, 15, 17, 51, 49, 123456789, time.Local)
	tx := timefy.New(t)           // or timefy.With(t)
	fmt.Println(tx.EndOfMonth()) // 2025-01-31 23:59:59.999999999 +0700 +07
}
```

> Calculating the time based on configuration

```go
package main

import (
	"fmt"
	"time"

	"github.com/sivaosorg/timefy"
)

func main() {
	// time based on configuration
	rule := timefy.NewRule().WithLocationRFC(timefy.DefaultTimezoneVietnam)
	t = time.Date(2025, 01, 15, 17, 51, 49, 123456789, time.Now().Location())
	fmt.Println(t)                              // 2025-01-15 17:51:49.123456789 +0700 +07
	fmt.Println(rule.With(t).BeginningOfWeek()) // 2025-01-12 00:00:00 +0700 +07
	v, _ := rule.Parse("2025-11-12 22:14:01")
	fmt.Println(v) // 2025-11-12 22:14:01 +0700 +07
}
```

> Monday/ Sunday

```go
package main

import (
	"fmt"
	"time"

	"github.com/sivaosorg/timefy"
)

func main() {
	// Monday / Sunday
	t := time.Date(2025, 01, 15, 17, 51, 49, 123456789, time.Now().Location())
	fmt.Println(timefy.Monday())                   // 2024-10-21 00:00:00 +0700 +07
	fmt.Println(timefy.Monday("15:35:34"))         // 2024-10-21 15:35:34 +0700 +07
	fmt.Println(timefy.Sunday())                   // 2024-10-27 00:00:00 +0700 +07
	fmt.Println(timefy.Sunday("16:45:34"))         // 2024-10-27 16:45:34 +0700 +07
	fmt.Println(timefy.EndOfSunday())              // 2024-10-27 23:59:59.999999999 +0700 +07
	fmt.Println(timefy.With(t).Monday())           // 2025-01-13 00:00:00 +0700 +07
	fmt.Println(timefy.With(t).Monday("15:35:34")) // 2025-01-13 15:35:34 +0700 +07
	fmt.Println(timefy.With(t).Sunday())           // 2025-01-19 00:00:00 +0700 +07
	fmt.Println(timefy.With(t).Sunday("16:45:34")) // 2025-01-19 16:45:34 +0700 +07
	fmt.Println(timefy.With(t).EndOfSunday())      // 2025-01-19 23:59:59.999999999 +0700 +07
}
```

> Parse String to Time

```go
package main

import (
	"fmt"

	"github.com/sivaosorg/timefy"
)

func main() {
	// String to Time
	t, _ := timefy.Parse("2025")
	fmt.Println(t) // 2025-01-01 00:00:00 +0700 +07
	t, _ = timefy.Parse("2025-02")
	fmt.Println(t) // 2025-02-01 00:00:00 +0700 +07
	t, _ = timefy.Parse("2025-02-15")
	fmt.Println(t) // 2025-02-15 00:00:00 +0700 +07
	t, _ = timefy.Parse("11-15")
	fmt.Println(t) // 2024-11-15 00:00:00 +0700 +07
	t, _ = timefy.Parse("13:15")
	fmt.Println(t) // 2024-10-25 13:15:00 +0700 +07
	t, _ = timefy.Parse("13:15:45")
	fmt.Println(t) // 2024-10-25 13:15:45 +0700 +07
	t, _ = timefy.Parse("23")
	fmt.Println(t) // 2024-10-25 23:00:00 +0700 +07
	// MustParse must parse string to time or it will panic
	t = timefy.MustParse("11")
	fmt.Println(t) // 2024-10-25 11:00:00 +0700 +07
	t = timefy.MustParse("99:99") // panic: can't parse string as time: 99:99
	// fmt.Println(t)

	// Extend timefy to support more formats is quite easy, just update timefy.TimeFormats with other time layouts, e.g:
	timefy.TimeFormats = append(timefy.TimeFormats, "02 Jan 2006 15:04")
	fmt.Println(timefy.TimeFormats)
}
```

### Contributing

To contribute to project, follow these steps:

1. Clone the repository:

   ```bash
   git clone --depth 1 https://github.com/sivaosorg/timefy.git
   ```

2. Navigate to the project directory:

   ```bash
   cd timefy
   ```

3. Prepare the project environment:
   ```bash
   go mod tidy
   ```
