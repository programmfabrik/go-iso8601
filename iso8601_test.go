package iso8601

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type dateCase struct {
	input string
	isErr bool
	width Width
	from  string
	to    string
}

func (dc dateCase) test() {

}

func TestCases(t *testing.T) {
	for _, tCase := range []dateCase{
		{"-200000", false, WIDTH_YEAR, "-200000-01-01T00:00:00Z", "-200000-12-31T23:59:59Z"},
		{"+200000", false, WIDTH_YEAR, "200000-01-01T00:00:00Z", "200000-12-31T23:59:59Z"},
		{"-200000-12", false, WIDTH_YEAR_MONTH, "-200000-12-01T00:00:00Z", "-200000-12-31T23:59:59Z"},
		{"-200000-12-12", false, WIDTH_DATE, "-200000-12-12T00:00:00Z", "-200000-12-12T23:59:59Z"},
		// {"-0000", true, year, "", ""},
		{"+0000", false, WIDTH_YEAR, "0000-01-01T00:00:00Z", "0000-12-31T23:59:59Z"},
		{"+0400", false, WIDTH_YEAR, "0400-01-01T00:00:00Z", "0400-12-31T23:59:59Z"},
		{"0400", false, WIDTH_YEAR, "0400-01-01T00:00:00Z", "0400-12-31T23:59:59Z"},
		{"-0400", false, WIDTH_YEAR, "-0400-01-01T00:00:00Z", "-0400-12-31T23:59:59Z"},
		{"-4", true, WIDTH_YEAR, "", ""},
		{"43", true, WIDTH_YEAR, "", ""},
		{"400", true, WIDTH_YEAR, "", ""},
		{"2000", false, WIDTH_YEAR, "2000-01-01T00:00:00Z", "2000-12-31T23:59:59Z"},
		// {"2000+01:00", false, WIDTH_YEAR, "2000-01-01T00:00:00+01:00", "2000-12-31T23:59:59+01:00"},
		{"2000-10", false, WIDTH_YEAR_MONTH, "2000-10-01T00:00:00Z", "2000-10-31T23:59:59Z"},
		{"1999-10-12", false, WIDTH_DATE, "1999-10-12T00:00:00Z", "1999-10-12T23:59:59Z"},
		{"1900-02-29", true, WIDTH_DATE, "", ""}, // No leap year
		{"2001-02-29", true, WIDTH_DATE, "", ""}, // No leap year
		{"2000-02-29", false, WIDTH_DATE, "2000-02-29T00:00:00Z", "2000-02-29T23:59:59Z"},
		// Allow leap year < 0
		{"-1972-02-29", false, WIDTH_DATE, "-1972-02-29T00:00:00Z", "-1972-02-29T23:59:59Z"},
		{"1972:06:27", false, WIDTH_DATE, "1972-06-27T00:00:00Z", "1972-06-27T23:59:59Z"},
		{"1899-12-31T23:59:59Z", false, WIDTH_DATE_TIME_SEC_TZ, "1899-12-31T23:59:59Z", "1899-12-31T23:59:59Z"},
		{"1399-12-31T23:59:59.000Z", false, WIDTH_DATE_TIME_SEC_TZ, "1399-12-31T23:59:59Z", "1399-12-31T23:59:59Z"},
		{"1999-10-12T12:20", false, WIDTH_DATE_TIME, "1999-10-12T12:20:00Z", "1999-10-12T12:20:59Z"},
		{"1999-10-12T12:20:10", false, WIDTH_DATE_TIME_SEC, "1999-10-12T12:20:10Z", "1999-10-12T12:20:10Z"},
		{"1999-10-12 12:20:10", false, WIDTH_DATE_TIME_SEC, "1999-10-12T12:20:10Z", "1999-10-12T12:20:10Z"},
		{"1999:10:12 12:20:10", false, WIDTH_DATE_TIME_SEC, "1999-10-12T12:20:10Z", "1999-10-12T12:20:10Z"},
		{"1999-10-12T12:20:10.12345", false, WIDTH_DATE_TIME_SEC, "1999-10-12T12:20:10Z", "1999-10-12T12:20:10Z"},
		{"1999-10-12T12:20:10+0200", false, WIDTH_DATE_TIME_SEC_TZ, "1999-10-12T12:20:10+02:00", "1999-10-12T12:20:10+02:00"},
		{"1999-10-12T12:20:10.123+0200", false, WIDTH_DATE_TIME_SEC_TZ, "1999-10-12T12:20:10+02:00", "1999-10-12T12:20:10+02:00"},
		{"1999-10-12T12:20:10+0200", false, WIDTH_DATE_TIME_SEC_TZ, "1999-10-12T12:20:10+02:00", "1999-10-12T12:20:10+02:00"},
		{"1999-10-12T12:20:10.123+02:00", false, WIDTH_DATE_TIME_SEC_TZ, "1999-10-12T12:20:10+02:00", "1999-10-12T12:20:10+02:00"},
		{"1999-10-12T12:20:10.123+02:00:10", false, WIDTH_DATE_TIME_SEC_TZS, "1999-10-12T12:20:10+02:00", "1999-10-12T12:20:10+02:00"},
		{"1999-10-12T12:20:10Z", false, WIDTH_DATE_TIME_SEC_TZ, "1999-10-12T12:20:10Z", "1999-10-12T12:20:10Z"},
		{"1999-10-12T12:20:10Z", false, WIDTH_DATE_TIME_SEC_TZ, "1999-10-12T12:20:10Z", "1999-10-12T12:20:10Z"},
		{"-0111-07-01T00:00:00+00:53:28", false, WIDTH_DATE_TIME_SEC_TZS, "-0111-07-01T00:00:00+00:53", "-0111-07-01T00:00:00+00:53"},
		{"2015-01-06 12:43:53+01", false, WIDTH_DATE_TIME_SEC_TZ, "2015-01-06T12:43:53+01:00", "2015-01-06T12:43:53+01:00"},
		{"2023-08-29 14:59:12 +01:00", false, WIDTH_DATE_TIME_SEC_TZ, "2023-08-29T14:59:12+01:00", "2023-08-29T14:59:12+01:00"},
		{"2017:05:22 16:07:28Z", false, WIDTH_DATE_TIME_SEC_TZ, "2017-05-22T16:07:28Z", "2017-05-22T16:07:28Z"},
		{"2023-08-29 14:59", false, WIDTH_DATE_TIME, "2023-08-29T14:59:00Z", "2023-08-29T14:59:59Z"},
	} {
		tm, err := Parse(tCase.input)
		if tCase.isErr {
			assert.Error(t, err, "%q", tCase.input)
			continue
		}
		if !assert.NoError(t, err, "%q", tCase.input) {
			return
		}
		if !assert.Equal(t, tCase.width, tm.Width, "width don't match %q", tCase.input) {
			return
		}

		from := tm.From.Format(time.RFC3339)
		to := tm.To.Format(time.RFC3339)
		if !assert.Equal(t, tCase.from, from, "%q", tCase.input) {
			return
		}
		if !assert.Equal(t, tCase.to, to, "%q", tCase.input) {
			return
		}
	}
}

func TestChangeWidth(t *testing.T) {
	pd := "2022-03-08T12:20:01"
	tm, err := Parse(pd)
	assert.NoError(t, err)
	assert.Equal(t, "2022-03-08T12:20:01", tm.FromString())
	tm.Width = WIDTH_DATE
	assert.Equal(t, "2022-03-08", tm.FromString())
}

func TestParseInLocation(t *testing.T) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if !assert.NoError(t, err) {
		return
	}
	_ = loc
	tm, err := Parse("2024-01-01")
	if !assert.NoError(t, err) {
		return
	}
	tm2 := tm.SetLocation(loc)
	if !assert.Equal(t, "2024-01-01 00:00:00 -0800 PST", tm2.From.String()) {
		return
	}
	if !assert.Equal(t, "2024-01-01 23:59:59 -0800 PST", tm2.To.String()) {
		return
	}
	if !assert.Equal(t, "2024-01-01T00:00:00-08:00-2024-01-01T23:59:59-08:00", tm2.String()) {
		return
	}
}

func TestLimits(t *testing.T) {
	tm, err := Parse("-9223372036854775807")
	if !assert.Error(t, err) {
		// range error
		return
	}

	// Below lower limit
	tm, err = Parse("-292277022400-12-31 20:00:01")
	if !assert.Error(t, err) {
		return
	}

	// Lower limit
	tm, err = Parse("-292277022399-12-31 20:00:01")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, "-292277022399-12-31T20:00:01", tm.String()) {
		return
	}

	// Upper limit
	tm, err = Parse("292277022399-12-31 20:00:01")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, "292277022399-12-31T20:00:01", tm.String()) {
		return
	}

	// Over Upper limit
	tm, err = Parse("292277022400-12-31 20:00:01")
	if !assert.Error(t, err) {
		return
	}

	// range checks

	// y := -292277020000
	// for {
	// 	tmin := time.Date(y, 1, 1, 0, 0, 0, 0, time.UTC)
	// 	if tmin.Year() != y {
	// 		golib.Pln("first break %d year %d format %s rfc %s", y, tmin.Year(), tmin.Format("2006"), tmin.Format(time.RFC3339))
	// 		break
	// 	}
	// 	y--
	// }

	// y = 292277024000
	// for {
	// 	tmin := time.Date(y, 1, 1, 0, 0, 0, 0, time.UTC)
	// 	if tmin.Year() != y {
	// 		golib.Pln("last break %d year %d format %s rfc %s", y, tmin.Year(), tmin.Format("2006"), tmin.Format(time.RFC3339))
	// 		break
	// 	}
	// 	y++
	// }

}

func TestEmpty(t *testing.T) {
	_, err := Parse("")
	if !assert.Error(t, err) {
		return
	}
}

func TestAdjustTimeZone(t *testing.T) {
	for _, dates := range [][2]string{
		{"2015-11-25T00:09:06-18:00", "2015-11-25T00:09:06-18:00"},
		{"2015-11-25T00:09:06-24:00", "2015-11-26T00:09:06Z"},
		{"2015-11-25T00:09:06-23:00", "2015-11-26T00:09:06+01:00"},
	} {
		tm, err := time.Parse(time.RFC3339, dates[0])
		if !assert.NoError(t, err) {
			return
		}
		tm2 := adjustTimeZone(tm, -18, 18)
		if !assert.Equal(t, dates[1], tm2.Format(time.RFC3339)) {
			return
		}
	}
}
