package clock

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Time deals only with hours, minutes, and seconds.  As opposed to the native
// time.Time struct that deals also with dates & timezones.
type Time struct {
	hours, minutes, seconds int
}

var (
	// StartOfDayTime is the time referring to 00:00:00
	StartOfDayTime Time

	// EndOfDayTime is the time referring to 23:59:59
	EndOfDayTime Time

	// ErrInvalidTimeFormat indicates that the input string is in
	// an invalid format.
	ErrInvalidTimeFormat = errors.New("invalid time format")
)

func init() {
	StartOfDayTime = NewTime(0, 0, 0)
	EndOfDayTime = NewTime(23, 59, 59)
}

// NewTime returns a new Time object given hours, minutes, and seconds.
func NewTime(h, m, s int) Time {
	return Time{
		hours:   h,
		minutes: m,
		seconds: s,
	}
}

func loadTimeZone(timezone string) *time.Location {
	loc := time.UTC
	tzLoc, err := time.LoadLocation(timezone)
	if err == nil {
		loc = tzLoc
	}

	return loc
}

// ParseTime takes in a string of the format: hh:mm:ss
// and returns a parsed Time object.  If the string is not
// in a valid format ErrInvalidTimeFormat is returned.
func ParseTime(str string) (*Time, error) {
	split := strings.Split(str, ":")
	if len(split) != 3 {
		return nil, fmt.Errorf("string not in form hh:mm:ss - %w", ErrInvalidTimeFormat)
	}

	hours, err := strconv.Atoi(split[0])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, ErrInvalidTimeFormat)
	}

	minutes, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, ErrInvalidTimeFormat)
	}

	seconds, err := strconv.Atoi(split[2])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, ErrInvalidTimeFormat)
	}

	tm := NewTime(hours, minutes, seconds)

	return &tm, nil
}

// Now returns the current Time at the sepcified timezone.
// If an invalid timezone is given, then UTC is used.
func Now(timezone string) Time {
	loc := loadTimeZone(timezone)
	now := time.Now().In(loc)

	return NewTime(now.Hour(), now.Minute(), now.Second())
}

// Today converts the Time object into a time.Time at the current
// date given a timezone.
func (t *Time) Today(timezone string) time.Time {
	loc := loadTimeZone(timezone)

	hours, minutes, seconds := t.HoursMinutesSeconds()
	now := time.Now().In(loc)
	return time.Date(now.Year(), now.Month(), now.Day(), hours, minutes, seconds, 0, loc)
}

func digitString(n int) string {
	str := strconv.Itoa(n)
	if n < 10 {
		str = fmt.Sprintf("0%s", str)
	}

	return str
}

// String returns the string representation of Time: hh:mm:ss
func (t *Time) String() string {
	return fmt.Sprintf(
		"%s:%s:%s",
		digitString(t.hours),
		digitString(t.minutes),
		digitString(t.seconds),
	)
}

// Value implements the sql.Valuer interface so that Time can be used
// in conjunction with the time type in databases.
func (t *Time) Value() (driver.Value, error) {
	return t.String(), nil
}

// Scan implements the sql.Scanner interface.
func (t *Time) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	bb, ok := src.([]byte)
	if !ok {
		return errors.New("failed to parse util.Time from sql driver")
	}

	str := string(bb)
	parsedTime, err := ParseTime(str)
	if err != nil {
		return err
	}

	*t = *parsedTime

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Time) UnmarshalJSON(data []byte) error {
	str := string(data)
	tt, err := ParseTime(strings.ReplaceAll(str, `"`, ""))
	if err != nil {
		return err
	}

	*t = *tt

	return nil
}

// dateTime is an internal method for converting the Time to
// an arbitrary time.Time.  This is used internally for computing addition
// and subtraction on Time.
func (t Time) dateTime() time.Time {
	return time.Date(2000, 1, 1, t.hours, t.minutes, t.seconds, 0, time.UTC)
}

// Add increments the Time by the given input duration.
func (t Time) Add(d time.Duration) Time {
	newDateTime := t.dateTime().Add(d)

	return NewTime(newDateTime.Hour(), newDateTime.Minute(), newDateTime.Second())
}

// Sub decrements the Time by the given input duration.
func (t Time) Sub(d time.Duration) Time {
	return t.Add(-d)
}

// HoursMinutesSeconds returns the given hours, minutes, and seconds within the Time object.
func (t Time) HoursMinutesSeconds() (int, int, int) {
	return t.hours, t.minutes, t.seconds
}

// TotalSeconds returns the total amount of seconds into the day that this Time object is.
func (t Time) TotalSeconds() int {
	h, m, s := t.HoursMinutesSeconds()
	return h*60*60 + m*60 + s
}

// After returns true fo the Time object occurs after the input Time.
func (t Time) After(comparison Time) bool {
	return t.TotalSeconds() > comparison.TotalSeconds()
}

// Before returns true fo the Time object occurs before the input Time.
func (t Time) Before(comparison Time) bool {
	return !t.After(comparison)
}

// DurationBetween returns the duration between the two times.  If start occurs
// after end, then the returned duration assumes that end refers to the following day.
func DurationBetween(start Time, end Time) time.Duration {
	if start.After(end) {
		return time.Duration(
			EndOfDayTime.TotalSeconds()-start.TotalSeconds()+end.TotalSeconds(),
		) * time.Second
	}

	return time.Duration(end.TotalSeconds()-start.TotalSeconds()) * time.Second
}

// Within returns true if the Time occurs within the start and end range.  If start occurs
// after end, then the returned duration assumes that end refers to the following day.
func (t Time) Within(start Time, end Time) bool {
	if start.After(end) {
		return (t.After(StartOfDayTime) && t.Before(end)) ||
			(t.After(start) && t.Before(EndOfDayTime))
	}

	return t.After(start) && t.Before(end)
}

// TimePointer returns a pointer reference to the input Time object.
func TimePointer(t time.Time) *time.Time {
	return &t
}
