// iso8601 is a ISO 8601 compliant date parser.
// The Parse function returns a Go STDLIB compliant time.Time
package iso8601

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Width int

const (
	WIDTH_YEAR Width = iota
	WIDTH_YEAR_MONTH
	WIDTH_DATE
	WIDTH_DATE_TIME
	WIDTH_DATE_TIME_SEC
	WIDTH_DATE_TIME_SEC_TZ
	WIDTH_DATE_TIME_SEC_TZS
)

const (
	YEAR              = "2006"
	YEAR_MONTH        = "2006-01"
	DATE              = "2006-01-02"
	DATE_TIME         = "2006-01-02T15:04"
	DATE_TIME_SEC     = "2006-01-02T15:04:05"
	DATE_TIME_SEC_TZ  = "2006-01-02T15:04:05Z07:00"
	DATE_TIME_SEC_TZS = "2006-01-02T15:04:05Z07:00:00"
)

var formats = map[Width]string{
	WIDTH_YEAR:              YEAR,
	WIDTH_YEAR_MONTH:        YEAR_MONTH,
	WIDTH_DATE:              DATE,
	WIDTH_DATE_TIME:         DATE_TIME,
	WIDTH_DATE_TIME_SEC:     DATE_TIME_SEC,
	WIDTH_DATE_TIME_SEC_TZ:  DATE_TIME_SEC_TZ,
	WIDTH_DATE_TIME_SEC_TZS: DATE_TIME_SEC_TZS,
}

var parseLayouts = []layout{
	{WIDTH_YEAR, YEAR},
	{WIDTH_YEAR_MONTH, YEAR_MONTH},
	{WIDTH_DATE, DATE},
	{WIDTH_DATE_TIME, DATE_TIME},
	{WIDTH_DATE_TIME_SEC, DATE_TIME_SEC},
	{WIDTH_DATE_TIME_SEC, "2006-01-02 15:04:05.999999999"},
	{WIDTH_DATE_TIME_SEC, "2006-01-02T15:04:05.999999999"},
	{WIDTH_DATE_TIME_SEC_TZS, DATE_TIME_SEC_TZS},
	{WIDTH_DATE_TIME_SEC_TZS, "2006-01-02T15:04:05.999999999Z070000"},
	{WIDTH_DATE_TIME_SEC_TZS, "2006-01-02T15:04:05.999999999-070000"},
	{WIDTH_DATE_TIME_SEC_TZS, "2006-01-02T15:04:05.999999999Z07:00:00"},
	{WIDTH_DATE_TIME_SEC_TZS, "2006-01-02T15:04:05.999999999-07:00:00"},
	{WIDTH_DATE_TIME_SEC_TZ, "2006-01-02T15:04:05.999999999Z070000"},
	{WIDTH_DATE_TIME_SEC_TZ, "2006-01-02T15:04:05.999999999-070000"},
	{WIDTH_DATE_TIME_SEC_TZ, "2006-01-02T15:04:05.999999999Z07:00"},
	{WIDTH_DATE_TIME_SEC_TZ, "2006-01-02T15:04:05.999999999-07:00"},
	{WIDTH_DATE_TIME_SEC_TZ, DATE_TIME_SEC_TZ},
	{WIDTH_DATE_TIME_SEC_TZ, "2006-01-02T15:04:05.999999999Z0700"},
	{WIDTH_DATE_TIME_SEC_TZ, "2006-01-02T15:04:05.999999999-0700"},
	{WIDTH_DATE_TIME_SEC_TZ, "2006-01-02T15:04:05.999999999Z07:00"},
	{WIDTH_DATE_TIME_SEC_TZ, "2006-01-02T15:04:05.999999999-07:00"},
	{WIDTH_DATE_TIME_SEC_TZ, "2006-01-02T15:04:05.999999999Z07"},
	{WIDTH_DATE_TIME_SEC_TZ, "2006-01-02T15:04:05.999999999-07"},
	{WIDTH_DATE_TIME_SEC_TZ, "2006-01-02T15:04:05.999999999Z07"},
	{WIDTH_DATE_TIME_SEC_TZ, "2006-01-02T15:04:05.999999999-07"},
	// - and " "
	{WIDTH_DATE_TIME, "2006-01-02 15:04"},
	{WIDTH_DATE_TIME_SEC_TZ, "2006-01-02 15:04:05.999999999Z07"},
	{WIDTH_DATE_TIME_SEC_TZ, "2006-01-02 15:04:05.999999999-07"},
	{WIDTH_DATE_TIME_SEC_TZ, "2006-01-02 15:04:05.999999999 Z07"},
	{WIDTH_DATE_TIME_SEC_TZ, "2006-01-02 15:04:05.999999999 -07"},
	{WIDTH_DATE_TIME_SEC_TZ, "2006-01-02 15:04:05.999999999 -0700"},
	{WIDTH_DATE_TIME_SEC_TZ, "2006-01-02 15:04:05.999999999 -07:00"},
	// refrain from handling abbreviated timezone, e.g. "CET" is recognized but "JST" is not, but silently accepted
	// so sorry Go, this is unusable
	// {WIDTH_DATE_TIME_SEC_TZ, "2006-01-02 15:04:05.999999999 MST"},

	// : separator
	{WIDTH_DATE, "2006:01:02"},                                      // format used by exiftool
	{WIDTH_DATE_TIME_SEC, "2006:01:02 15:04:05.999999999"},          // format used by exiftool
	{WIDTH_DATE_TIME_SEC_TZ, "2006:01:02 15:04:05.999999999-07:00"}, // format used by exiftool
	{WIDTH_DATE_TIME_SEC_TZ, "2006:01:02 15:04:05.999999999-07"},    // format used by exiftool
	{WIDTH_DATE_TIME_SEC_TZ, "2006:01:02 15:04:05.999999999Z07"},    // format used by exiftool
	{WIDTH_DATE_TIME_SEC_TZ, time.RFC3339},
}

type Time struct {
	From  time.Time
	To    time.Time
	Width Width
}

func init() {
	// sort layouts by width. start with the lowest. within the same width,
	// the longer layout string comes first
	slices.SortFunc(parseLayouts, func(a, b layout) int {
		switch {
		case a.width < b.width:
			return -1
		case a.width > b.width:
			return 1
		default:
			la := len(a.layout)
			lb := len(b.layout)
			switch {
			case la < lb:
				return 1
			case lb > la:
				return -1
			default:
				return 0
			}
		}
	})
}

// AdjustTimeZone moves the timezone to with from and to offsets. This is a helper func
// to be able to format Java ZoneOffset compatible dates. That function limits the range from
// -18:00 to 18:00, whereas Go accepts -24:59 to 24:59.
func (t *Time) AdjustTimeZone(fromH, toH int) (t2 *Time) {
	if t == nil {
		return t
	}
	t.From = adjustTimeZone(t.From, fromH, toH)
	t.To = adjustTimeZone(t.To, fromH, toH)
	return t
}

func adjustTimeZone(t time.Time, fromH, toH int) (t2 time.Time) {
	_, offset := t.Zone()
	from := fromH * 3600
	to := toH * 3600
	if offset >= from && offset <= to {
		return t
	}
	var fz *time.Location
	if offset < from {
		fz = time.FixedZone("", offset+(24*3600)) // -23 > +1
	} else if offset > to {
		fz = time.FixedZone("", offset-(24*3600)) // 23 > -1
	}
	return t.UTC().In(fz)
}

// String formats the time returning the same value
// which was parsed, if the width wasn't changed
func (t Time) String() string {
	from := t.FromString()
	to := t.ToString()
	if from == to {
		return from
	}
	return from + "-" + to
}

func (t Time) FromString() string {
	return t.From.Format(formats[t.Width])
}

func (t Time) ToString() string {
	return t.To.Format(formats[t.Width])
}

// Format returns the format of the stored time. The format is reflecting the
// width of the stored value.
func (t Time) Format() string {
	return formats[t.Width]
}

// HasTime returns true if the format contains time information
func (t Time) HasTime() bool {
	switch t.Width {
	case WIDTH_YEAR, WIDTH_YEAR_MONTH, WIDTH_DATE:
		return false
	}
	return true
}

type layout struct {
	width  Width
	layout string
}

// daysIn returns the number of days in a month for a given year.
func daysIn(m time.Month, yearI int) int {
	// This is equivalent to time.daysIn(m, year).
	return time.Date(yearI, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// timeToLoc strips of time location information and sets
// the location to loc. This changes the timestamp of t.
func timeToLoc(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
}

// SetLocation overwrites the timezone in t to match loc. This changes the width
// of t to WIDTH_DATE_TIME_SEC_TZS. If loc is nil, an unmodified copy of t is
// returned.
func (t Time) SetLocation(loc *time.Location) (t2 Time) {
	t2 = t
	if loc == nil {
		return t2
	}
	if !t2.From.IsZero() {
		t2.From = timeToLoc(t2.From, loc)
	}
	if !t2.To.IsZero() {
		t2.To = timeToLoc(t2.To, loc)
	}
	t2.Width = WIDTH_DATE_TIME_SEC_TZ
	return t2
}

// Parse parse ISO 8601 compatible layouts. In addition to standard YYYY representation
// this accepts an arbitrary number of digits for the yearI, as well as a "-" or "+" preceding it.
// For the year 1 BC, this expects +0000.
// For dates BC we always assume a leap yearI, so that if someone want to parse "02-29" we are
// not getting in her way.
func Parse(dOrig string) (t *Time, err error) {

	t = &Time{}

	d := dOrig

	// Parse the YEAR at the beginning of the date string to support
	// +/- before the years
	re := regexp.MustCompile(`^((-|\+|)[0-9]+)`) // ((:|-|+).+|$)`)
	yearS := re.FindStringSubmatch(d)
	if yearS == nil {
		return nil, fmt.Errorf("Unable to parse date %q", d)
	}

	yearI, err := strconv.Atoi(yearS[1])
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to parse %q", d)
	}

	// if yearI == 0 && yearS[2] != "+" {
	// 	return nil, fmt.Errorf("Unable to parse 1 BC: %q. This must be +0000.", d)
	// }

	// matched year "-2000" - length of the "-""
	if len(yearS[1])-len(yearS[2]) < 4 {
		return nil, fmt.Errorf("Unable to parse %q: Year needs to be at least 4 digits", d)
	}

	// We split of the YEAR manually and feed always years 2000 / 2001 to
	// time.Parse. This is to support years which are not supported well by the stdlib.

	correctedYear := false
	if yearI < 1000 || yearI > 9999 {
		isLeap := yearI <= 0 || (yearI%400 == 0 || yearI%4 == 0 && yearI%100 != 0)
		if isLeap {
			d = "2000" + d[len(yearS[1]):]
			yearI -= 2000
		} else {
			d = "2001" + d[len(yearS[1]):]
			yearI -= 2001
		}
		correctedYear = true
	}
	for _, layout := range parseLayouts {
		switch layout.width {
		case WIDTH_DATE_TIME_SEC, WIDTH_DATE_TIME_SEC_TZ, WIDTH_DATE_TIME_SEC_TZS:
			// some timezone date layouts have variable lengths
		default:
			if len(d) != len(layout.layout) {
				// Dates with fixed length, we only need to parse if the input length matches
				continue
			}
		}

		t2, err := time.Parse(layout.layout, d)
		if err != nil {
			// println(d, layout.layout, err.Error())
			continue
		}

		// golib.Pln("parsed %q %q %s", layout.layout, d, t2.String())

		t2 = t2.Truncate(time.Second)
		t.Width = layout.width

		if correctedYear {
			yearI += t2.Year()
		} else {
			yearI = t2.Year()
		}

		if yearI < -292277022399 || yearI > 292277022399 {
			// technically time.Time can work with years up to 292277026853, but
			// to keep the min & max identical we do it like this
			return nil, fmt.Errorf("year %d ouside of allowed range -292277022399 - 292277022399", yearI)
		}

		switch layout.width {
		case WIDTH_YEAR:
			t.From = time.Date(yearI, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
			t.To = time.Date(yearI, time.Month(12), 31, 23, 59, 59, 0, time.UTC)
		case WIDTH_YEAR_MONTH:
			t.From = time.Date(yearI, t2.Month(), 1, 0, 0, 0, 0, time.UTC)
			t.To = time.Date(yearI, t2.Month(), daysIn(t2.Month(), yearI), 23, 59, 59, 0, time.UTC)
		case WIDTH_DATE:
			t.From = time.Date(yearI, t2.Month(), t2.Day(), 0, 0, 0, 0, time.UTC)
			t.To = time.Date(yearI, t2.Month(), t2.Day(), 23, 59, 59, 0, time.UTC)
		case WIDTH_DATE_TIME:
			t.From = time.Date(yearI, t2.Month(), t2.Day(), t2.Hour(), t2.Minute(), 0, 0, time.UTC)
			t.To = time.Date(yearI, t2.Month(), t2.Day(), t2.Hour(), t2.Minute(), 59, 0, time.UTC)
		case WIDTH_DATE_TIME_SEC:
			t.From = time.Date(yearI, t2.Month(), t2.Day(), t2.Hour(), t2.Minute(), t2.Second(), 0, time.UTC)
			t.To = time.Date(yearI, t2.Month(), t2.Day(), t2.Hour(), t2.Minute(), t2.Second(), 0, time.UTC)
		case WIDTH_DATE_TIME_SEC_TZ, WIDTH_DATE_TIME_SEC_TZS:
			t.From = time.Date(yearI, t2.Month(), t2.Day(), t2.Hour(), t2.Minute(), t2.Second(), 0, t2.Location())
			t.To = time.Date(yearI, t2.Month(), t2.Day(), t2.Hour(), t2.Minute(), t2.Second(), 0, t2.Location())
		}

		return t, nil
	}

	return nil, fmt.Errorf("Unparsable date %q", dOrig)
}
