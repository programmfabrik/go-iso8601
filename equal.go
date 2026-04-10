package iso8601

// Equal returns true if the times from & to are equal at UTC time. We use the
// String() method for comparison.
func (t *Time) Equals(u *Time) bool {
	if t == nil && u == nil {
		return true
	}
	if t == nil || u == nil {
		return false
	}
	t.From = t.From.UTC()
	t.To = t.To.UTC()
	if t.String() != u.String() {
		return false
	}
	return true
}
