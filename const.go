package timefy

import (
	"regexp"
	"time"
)

type TimeRFC string
type TimeFormatRFC string
type ZoneRFC string

// WeekStartDay set week start day, default is sunday
var WeekStartDay = time.Sunday

// DefaultConfig default config
var DefaultConfig *Config

const (
	// Time in format 15:04:05, e.g., 13:45:30
	TimeRFC01T150405 TimeRFC = "15:04:05"

	// Time in format 15.04.05, e.g., 13.45.30
	TimeRFC02D150405 TimeRFC = "15.04.05"

	// Time in format 15-04-05, e.g., 13-45-30
	TimeRFC03H150405 TimeRFC = "15-04-05"

	// Time in format 15/04/05, e.g., 13/45/30
	TimeRFC04M150405 TimeRFC = "15/04/05"

	// Time in format 15 04 05, e.g., 13 45 30
	TimeRFC05S150405 TimeRFC = "15 04 05"

	// Time in format 150405, e.g., 134530
	TimeRFC06T150405 TimeRFC = "150405"

	// Time in format 15:04, e.g., 13:45
	TimeRFC07D150405 TimeRFC = "15:04"

	// Time in format 15.04, e.g., 13.45
	TimeRFC08H150405 TimeRFC = "15.04"

	// Time in format 15-04, e.g., 13-45
	TimeRFC09M150405 TimeRFC = "15-04"

	// Time in format 15/04, e.g., 13/45
	TimeRFC10S150405 TimeRFC = "15/04"

	// Time in format 1504, e.g., 1345
	TimeRFC11T150405 TimeRFC = "1504"
)

const (
	// Time in format 2006-01-02T15:04:05.999999, e.g., 2023-08-15T13:45:30.123456
	TimeFormat20060102T150405999999 TimeFormatRFC = "2006-01-02T15:04:05.999999"

	// Time in format 2006-01-02T15:04:05, e.g., 2023-08-15T13:45:30
	TimeFormat20060102T150405 TimeFormatRFC = "2006-01-02T15:04:05"

	// Time in format 2006-01-02 15:04:05, e.g., 2023-08-15 13:45:30
	TimeFormat20060102150405 TimeFormatRFC = "2006-01-02 15:04:05"

	// Time in format 02-01-2006 15:04:05, e.g., 15-08-2023 13:45:30
	TimeFormat02012006150405 TimeFormatRFC = "02-01-2006 15:04:05"

	// Time in format 02/01/2006 15:04:05, e.g., 15/08/2023 13:45:30
	TimeFormatRFC0102012006150405 TimeFormatRFC = "02/01/2006 15:04:05"

	// Time in format 2006-01-02 15:04:05.999999, e.g., 2023-08-15 13:45:30.123456
	TimeFormat20060102150405999999 TimeFormatRFC = "2006-01-02 15:04:05.999999"

	// Time in format 2006-01-02 15:04:05.999999999 -07:00, e.g., 2023-08-15 13:45:30.123456789 -07:00
	TimeFormat20060102150405999999RFC3339 TimeFormatRFC = "2006-01-02 15:04:05.999999999 -07:00"

	// Time in format 2006-01-02, e.g., 2023-08-15
	TimeFormat20060102 TimeFormatRFC = "2006-01-02"

	// Time in format 02/01/2006, e.g., 15/08/2023
	TimeFormatRFC0102012006 TimeFormatRFC = "02/01/2006"

	// Time in format 2006-01-02 15:04, e.g., 2023-08-15 13:45
	TimeFormat200601021504 TimeFormatRFC = "2006-01-02 15:04"

	// Time in format 2006-01-02 15, e.g., 2023-08-15 13
	TimeFormat2006010215 TimeFormatRFC = "2006-01-02 15"

	// Time in format 2006-01, e.g., 2023-08
	TimeFormat200601 TimeFormatRFC = "2006-01"

	// Time in format 02-01-2006, e.g., 15-08-2023
	TimeFormat02012006 TimeFormatRFC = "02-01-2006"

	// Time in format 01-02-2006, e.g., 08-15-2023
	TimeFormat01022006 TimeFormatRFC = "01-02-2006"

	// Time in format 2006-01-02 15:04:05 -07:00, e.g., 2023-08-15 13:45:30 -07:00
	TimeFormat20060102150405Z0700 TimeFormatRFC = "2006-01-02 15:04:05 -07:00"

	// Time in format 2006-01-02 15:04:05 -07:00:00, e.g., 2023-08-15 13:45:30 -07:00:00
	TimeFormat20060102150405Z070000 TimeFormatRFC = "2006-01-02 15:04:05 -07:00:00"

	// Time in format 2006-01-02T15:04:05-07:00, e.g., 2023-08-15T13:45:30-07:00
	TimeFormat20060102150405Z0700RFC3339 TimeFormatRFC = "2006-01-02T15:04:05-07:00"

	// Time in format 2006-01-02T15:04:05-07:00:00, e.g., 2023-08-15T13:45:30-07:00:00
	TimeFormat20060102150405Z070000RFC3339 TimeFormatRFC = "2006-01-02T15:04:05-07:00:00"

	// Time in format 2006-01-02 15:04:05 -07, e.g., 2023-08-15 13:45:30 -07
	TimeFormat20060102150405Z07 TimeFormatRFC = "2006-01-02 15:04:05 -07"

	// Time in format 2006-01-02T15:04:05-07, e.g., 2023-08-15T13:45:30-07
	TimeFormat20060102150405Z07RFC3339 TimeFormatRFC = "2006-01-02T15:04:05-07"

	// Time in format Mon, 02 Jan 2006 15:04:05 -0700, e.g., Tue, 15 Aug 2023 13:45:30 -0700
	TimeFormat20060102150405Z0700RFC1123 TimeFormatRFC = "Mon, 02 Jan 2006 15:04:05 -0700"

	// Time in format Mon, 02 Jan 2006 15:04:05 -07:00:00, e.g., Tue, 15 Aug 2023 13:45:30 -07:00:00
	TimeFormat20060102150405Z070000RFC1123 TimeFormatRFC = "Mon, 02 Jan 2006 15:04:05 -07:00:00"

	// Time in format Mon, 02 Jan 2006 15:04:05 -07, e.g., Tue, 15 Aug 2023 13:45:30 -07
	TimeFormat20060102150405Z07RFC1123 TimeFormatRFC = "Mon, 02 Jan 2006 15:04:05 -07"

	// Time in format 2006-01-02 15:04:05 UTC-07, e.g., 2023-08-15 13:45:30 UTC-07
	TimeFormat20060102150405Z07UTC TimeFormatRFC = "2006-01-02 15:04:05 UTC-07"

	// Time in format 2006-01-02 15:04:05 UTC-07:00, e.g., 2023-08-15 13:45:30 UTC-07:00
	TimeFormat20060102150405Z0700UTC TimeFormatRFC = "2006-01-02 15:04:05 UTC-07:00"

	// Time in format 2006-01-02 15:04:05 UTC-07:00:00, e.g., 2023-08-15 13:45:30 UTC-07:00:00
	TimeFormat20060102150405Z070000UTC TimeFormatRFC = "2006-01-02 15:04:05 UTC-07:00:00"

	// Time in format 2006-01-02T15:04:05UTC-07, e.g., 2023-08-15T13:45:30UTC-07
	TimeFormat20060102150405Z07UTCRFC3339 TimeFormatRFC = "2006-01-02T15:04:05UTC-07"

	// Time in format 2006-01-02T15:04:05UTC-07:00, e.g., 2023-08-15T13:45:30UTC-07:00
	TimeFormat20060102150405Z0700UTCRFC3339 TimeFormatRFC = "2006-01-02T15:04:05UTC-07:00"

	// Time in format 2006-01-02T15:04:05UTC-07:00:00, e.g., 2023-08-15T13:45:30UTC-07:00:00
	TimeFormat20060102150405Z070000UTCRFC3339 TimeFormatRFC = "2006-01-02T15:04:05UTC-07:00:00"
)

// Timezone constants representing default timezones for specific regions.
const (
	// DefaultTimezoneVietnam is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Vietnam, which is "Asia/Ho_Chi_Minh".
	DefaultTimezoneVietnam ZoneRFC = "Asia/Ho_Chi_Minh"

	// DefaultTimezoneNewYork is a constant that holds the IANA Time Zone identifier
	// for the default timezone in New York, USA, which is "America/New_York".
	DefaultTimezoneNewYork ZoneRFC = "America/New_York"

	// DefaultTimezoneLondon is a constant that holds the IANA Time Zone identifier
	// for the default timezone in London, United Kingdom, which is "Europe/London".
	DefaultTimezoneLondon ZoneRFC = "Europe/London"

	// DefaultTimezoneTokyo is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Tokyo, Japan, which is "Asia/Tokyo".
	DefaultTimezoneTokyo ZoneRFC = "Asia/Tokyo"

	// DefaultTimezoneSydney is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Sydney, Australia, which is "Australia/Sydney".
	DefaultTimezoneSydney ZoneRFC = "Australia/Sydney"

	// DefaultTimezoneParis is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Paris, France, which is "Europe/Paris".
	DefaultTimezoneParis ZoneRFC = "Europe/Paris"

	// DefaultTimezoneMoscow is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Moscow, Russia, which is "Europe/Moscow".
	DefaultTimezoneMoscow ZoneRFC = "Europe/Moscow"

	// DefaultTimezoneLosAngeles is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Los Angeles, USA, which is "America/Los_Angeles".
	DefaultTimezoneLosAngeles ZoneRFC = "America/Los_Angeles"

	// DefaultTimezoneManila is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Manila, Philippines, which is "Asia/Manila".
	DefaultTimezoneManila ZoneRFC = "Asia/Manila"

	// DefaultTimezoneKualaLumpur is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Kuala Lumpur, Malaysia, which is "Asia/Kuala_Lumpur".
	DefaultTimezoneKualaLumpur ZoneRFC = "Asia/Kuala_Lumpur"

	// DefaultTimezoneJakarta is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Jakarta, Indonesia, which is "Asia/Jakarta".
	DefaultTimezoneJakarta ZoneRFC = "Asia/Jakarta"

	// DefaultTimezoneYangon is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Yangon, Myanmar, which is "Asia/Yangon".
	DefaultTimezoneYangon ZoneRFC = "Asia/Yangon"

	// DefaultTimezoneAuckland is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Auckland, New Zealand, which is "Pacific/Auckland".
	DefaultTimezoneAuckland ZoneRFC = "Pacific/Auckland"

	// DefaultTimezoneBangkok is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Bangkok, Thailand, which is "Asia/Bangkok".
	DefaultTimezoneBangkok ZoneRFC = "Asia/Bangkok"

	// DefaultTimezoneDelhi is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Delhi, India, which is "Asia/Kolkata".
	DefaultTimezoneDelhi ZoneRFC = "Asia/Kolkata"

	// DefaultTimezoneDubai is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Dubai, United Arab Emirates, which is "Asia/Dubai".
	DefaultTimezoneDubai ZoneRFC = "Asia/Dubai"

	// DefaultTimezoneCairo is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Cairo, Egypt, which is "Africa/Cairo".
	DefaultTimezoneCairo ZoneRFC = "Africa/Cairo"

	// DefaultTimezoneAthens is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Athens, Greece, which is "Europe/Athens".
	DefaultTimezoneAthens ZoneRFC = "Europe/Athens"

	// DefaultTimezoneRome is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Rome, Italy, which is "Europe/Rome".
	DefaultTimezoneRome ZoneRFC = "Europe/Rome"

	// DefaultTimezoneJohannesburg is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Johannesburg, South Africa, which is "Africa/Johannesburg".
	DefaultTimezoneJohannesburg ZoneRFC = "Africa/Johannesburg"

	// DefaultTimezoneStockholm is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Stockholm, Sweden, which is "Europe/Stockholm".
	DefaultTimezoneStockholm ZoneRFC = "Europe/Stockholm"

	// DefaultTimezoneOslo is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Oslo, Norway, which is "Europe/Oslo".
	DefaultTimezoneOslo ZoneRFC = "Europe/Oslo"

	// DefaultTimezoneHelsinki is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Helsinki, Finland, which is "Europe/Helsinki".
	DefaultTimezoneHelsinki ZoneRFC = "Europe/Helsinki"

	// DefaultTimezoneKiev is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Kiev, Ukraine, which is "Europe/Kiev".
	DefaultTimezoneKiev ZoneRFC = "Europe/Kiev"

	// DefaultTimezoneBeijing is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Beijing, China, which is "Asia/Shanghai".
	DefaultTimezoneBeijing ZoneRFC = "Asia/Shanghai"

	// DefaultTimezoneSingapore is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Singapore, which is "Asia/Singapore".
	DefaultTimezoneSingapore ZoneRFC = "Asia/Singapore"

	// DefaultTimezoneIslamabad is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Islamabad, Pakistan, which is "Asia/Karachi".
	DefaultTimezoneIslamabad ZoneRFC = "Asia/Karachi"

	// DefaultTimezoneColombo is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Colombo, Sri Lanka, which is "Asia/Colombo".
	DefaultTimezoneColombo ZoneRFC = "Asia/Colombo"

	// DefaultTimezoneDhaka is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Dhaka, Bangladesh, which is "Asia/Dhaka".
	DefaultTimezoneDhaka ZoneRFC = "Asia/Dhaka"

	// DefaultTimezoneKathmandu is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Kathmandu, Nepal, which is "Asia/Kathmandu".
	DefaultTimezoneKathmandu ZoneRFC = "Asia/Kathmandu"

	// DefaultTimezoneBrisbane is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Brisbane, Australia, which is "Australia/Brisbane".
	DefaultTimezoneBrisbane ZoneRFC = "Australia/Brisbane"

	// DefaultTimezoneWellington is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Wellington, New Zealand, which is "Pacific/Auckland".
	DefaultTimezoneWellington ZoneRFC = "Pacific/Auckland"

	// DefaultTimezonePortMoresby is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Port Moresby, Papua New Guinea, which is "Pacific/Port_Moresby".
	DefaultTimezonePortMoresby ZoneRFC = "Pacific/Port_Moresby"

	// DefaultTimezoneSuva is a constant that holds the IANA Time Zone identifier
	// for the default timezone in Suva, Fiji, which is "Pacific/Fiji".
	DefaultTimezoneSuva ZoneRFC = "Pacific/Fiji"
)

var (
	// ApplyTimeRegexp is a regular expression that matches various time formats such as:
	// 15:04:05, 15:04:05.000, 15:04:05.000000, 15, 2017-01-01 15:04, 2021-07-20T00:59:10Z,
	// 2021-07-20T00:59:10+08:00, 2021-07-20T00:00:10-07:00, etc.
	ApplyTimeRegexp = regexp.MustCompile(`(\s+|^\s*|T)\d{1,2}((:\d{1,2})*|((:\d{1,2}){2}\.(\d{3}|\d{6}|\d{9})))(\s*$|[Z+-])`)

	// OnlyTimeRegexp is a regular expression that matches time formats such as:
	// 15:04:05, 15, 15:04:05.000, 15:04:05.000000, etc.
	OnlyTimeRegexp = regexp.MustCompile(`^\s*\d{1,2}((:\d{1,2})*|((:\d{1,2}){2}\.(\d{3}|\d{6}|\d{9})))\s*$`)
)

var (
	// TimeFormats is a slice of strings that holds various time format patterns.
	// These patterns can be used to parse and format time values in different layouts.
	TimeFormats = []string{
		"2006",                   // Year only, e.g., 2023
		"2006-1",                 // Year and month, e.g., 2023-8
		"2006-1-2",               // Year, month, and day, e.g., 2023-8-15
		"2006-1-2 15",            // Year, month, day, and hour, e.g., 2023-8-15 13
		"2006-1-2 15:4",          // Year, month, day, hour, and minute, e.g., 2023-8-15 13:45
		"2006-1-2 15:4:5",        // Year, month, day, hour, minute, and second, e.g., 2023-8-15 13:45:30
		"1-2",                    // Month and day, e.g., 8-15
		"15:4:5",                 // Hour, minute, and second, e.g., 13:45:30
		"15:4",                   // Hour and minute, e.g., 13:45
		"15",                     // Hour only, e.g., 13
		"15:4:5 Jan 2, 2006 MST", // Full date and time with timezone, e.g., 13:45:30 Aug 15, 2023 MST
		"2006-01-02 15:04:05.999999999 -0700 MST", // Full date and time with nanoseconds and timezone, e.g., 2023-08-15 13:45:30.123456789 -0700 MST
		"2006-01-02T15:04:05Z0700",                // ISO 8601 format with timezone offset, e.g., 2023-08-15T13:45:30Z0700
		"2006-01-02T15:04:05Z07",                  // ISO 8601 format with timezone offset (short), e.g., 2023-08-15T13:45:30Z07
		"2006.1.2",                                // Year, month, and day with dots, e.g., 2023.8.15
		"2006.1.2 15:04:05",                       // Year, month, day, and time with dots, e.g., 2023.8.15 13:45:30
		"2006.01.02",                              // Year, month, and day with leading zeros, e.g., 2023.08.15
		"2006.01.02 15:04:05",                     // Year, month, day, and time with leading zeros, e.g., 2023.08.15 13:45:30
		"2006.01.02 15:04:05.999999999",           // Full date and time with nanoseconds and leading zeros, e.g., 2023.08.15 13:45:30.123456789
		"1/2/2006",                                // Month, day, and year with slashes, e.g., 8/15/2023
		"1/2/2006 15:4:5",                         // Month, day, year, and time with slashes, e.g., 8/15/2023 13:45:30
		"2006/01/02",                              // Year, month, and day with slashes, e.g., 2023/08/15
		"20060102",                                // Year, month, and day without separators, e.g., 20230815
		"2006/01/02 15:04:05",                     // Year, month, day, and time with slashes, e.g., 2023/08/15 13:45:30
		time.ANSIC,                                // ANSIC format, e.g., Mon Aug 15 13:45:30 2023
		time.UnixDate,                             // Unix date format, e.g., Mon Aug 15 13:45:30 MST 2023
		time.RubyDate,                             // Ruby date format, e.g., Mon Aug 15 13:45:30 +0700 2023
		time.RFC822,                               // RFC 822 format, e.g., 15 Aug 23 13:45 MST
		time.RFC822Z,                              // RFC 822 format with numeric timezone, e.g., 15 Aug 23 13:45 -0700
		time.RFC850,                               // RFC 850 format, e.g., Monday, 15-Aug-23 13:45:30 MST
		time.RFC1123,                              // RFC 1123 format, e.g., Mon, 15 Aug 2023 13:45:30 MST
		time.RFC1123Z,                             // RFC 1123 format with numeric timezone, e.g., Mon, 15 Aug 2023 13:45:30 -0700
		time.RFC3339,                              // RFC 3339 format, e.g., 2023-08-15T13:45:30Z
		time.RFC3339Nano,                          // RFC 3339 format with nanoseconds, e.g., 2023-08-15T13:45:30.123456789Z
		time.Kitchen,                              // Kitchen format, e.g., 1:45PM
		time.Stamp,                                // Stamp format, e.g., Aug 15 13:45:30
		time.StampMilli,                           // Stamp format with milliseconds, e.g., Aug 15 13:45:30.123
		time.StampMicro,                           // Stamp format with microseconds, e.g., Aug 15 13:45:30.123456
		time.StampNano,                            // Stamp format with nanoseconds, e.g., Aug 15 13:45:30.123456789
	}
)
