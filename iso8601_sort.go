// iso8601 is a ISO 8601 compliant date parser.
// The Parse function returns a Go STDLIB compliant time.Time
package iso8601

import (
	"fmt"
	"math"
)

type sortMode int

const (
	sort_from sortMode = iota
	sort_to
	sort_middle
)

const maxInt64 = int64(^uint64(0)>>1) / 2 // 4611686018427387903
const minInt64 = -maxInt64                // -4611686018427387903

// FormatSortString returns a sortable string
// from an int64
func FormatSortString(i int64) string {
	return fmt.Sprintf("%022d", i+maxInt64)
}

// sort returns a sortable string value for time. The mode can be "from", "to",
// "middle" to pick an exact time for the sorting. With nullsMin and t == nil a
// minimal value for sorting is returned, if nullsMin is false the maximal value
// is returned. sort used a typed int64 represenatation to generate the string.
func (t *Time) sort(mode sortMode, nullsMin bool) string {
	if t == nil {
		if nullsMin {
			return FormatSortString(minInt64)
		} else {
			return FormatSortString(maxInt64)
		}
	}
	switch mode {
	case sort_from:
		return FormatSortString(t.From.Unix())
	case sort_to:
		return FormatSortString(t.To.Unix())
	case sort_middle:
		from := t.From.Unix()
		to := t.To.Unix()
		// use the to - from method to avoid overflowing ints
		// if we are close to maxInt
		return FormatSortString(int64(float64(from) + math.Round(float64(to-from)/2)))
	default:
		panic("Unknown mode")
	}
}

func (t *Time) SortTo(nullsMin bool) string {
	return t.sort(sort_to, nullsMin)
}

func (t *Time) SortFrom(nullsMin bool) string {
	return t.sort(sort_from, nullsMin)
}

// Sort returns a string for sorting the middle
// of the time
func (t *Time) Sort(nullMin bool) string {
	return t.sort(sort_middle, nullMin)
}
