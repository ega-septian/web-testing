package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// jsonTimeLayout is "yyyy-mm-dd HH:MM:SS" — plain and timezone-free in API
// responses, instead of Go's default RFC3339 (with offset and fractional seconds).
const jsonTimeLayout = "2006-01-02 15:04:05"

// JSONTime wraps time.Time purely to control its JSON representation. It
// still scans from Postgres exactly like time.Time (via sql.Scanner) — this
// only changes what goes out over the API, not how it's stored.
type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).Format(jsonTimeLayout) + `"`), nil
}

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		return nil
	}
	parsed, err := time.Parse(jsonTimeLayout, s)
	if err != nil {
		return err
	}
	*t = JSONTime(parsed)
	return nil
}

// Scan implements sql.Scanner so pgx can populate this field directly from a
// timestamptz column, the same way it would a plain time.Time.
func (t *JSONTime) Scan(value any) error {
	switch v := value.(type) {
	case time.Time:
		*t = JSONTime(v)
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into JSONTime", value)
	}
}

// Value implements driver.Valuer for symmetry, in case a JSONTime is ever
// used as a query parameter rather than just a scanned result.
func (t JSONTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}
