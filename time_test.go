package clock

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithin(t *testing.T) {
	t1 := NewTime(13, 01, 00)

	end := t1.Add(time.Hour)
	start := t1.Sub(time.Hour)

	assert.True(t, t1.Within(start, end))
	assert.False(t, t1.Within(end, start))

	start = t1.Add(2 * time.Hour)
	end = t1.Add(1 * time.Hour)
	assert.True(t, t1.Within(start, end))
	assert.False(t, t1.Within(end, start))
}

func TestJSON(t *testing.T) {
	rawJSON := `{
			"time": "10:11:12"
	}`

	var body struct {
		Time Time `json:"time"`
	}
	err := json.Unmarshal([]byte(rawJSON), &body)
	assert.Nil(t, err)

	marshaledBackRaw, err := json.Marshal(body)
	assert.Nil(t, err)
	assert.Equal(t, `{"time":"10:11:12"}`, string(marshaledBackRaw))
}

func TestDurationBetween(t *testing.T) {
	t1 := NewTime(00, 00, 00)

	start := t1.Sub(time.Hour)
	end := t1.Add(time.Hour)

	duration := DurationBetween(start, end)
	assert.Equal(t, 1*time.Hour+59*time.Minute+59*time.Second, duration)

	end = t1.Sub(time.Hour)
	start = t1.Add(time.Hour)

	duration = DurationBetween(start, end)
	assert.Equal(t, 22*time.Hour, duration)
}

func TestToday(t *testing.T) {
	timeToTest := NewTime(8, 30, 0)
	today := timeToTest.Today("US/Hawaii")
	hawaii, err := time.LoadLocation("US/Hawaii")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, time.Now().In(hawaii).Day(), today.Day())
	assert.Equal(t, 8, today.Hour())
	assert.Equal(t, 30, today.Minute())
	assert.Equal(t, 0, today.Second())
	assert.Equal(t, 18, today.UTC().Hour())
}

func TestString(t *testing.T) {
	tm := NewTime(8, 8, 8)
	assert.Equal(t, "08:08:08", tm.String())

	tm = NewTime(18, 8, 8)
	assert.Equal(t, "18:08:08", tm.String())

	tm = NewTime(12, 12, 12)
	assert.Equal(t, "12:12:12", tm.String())
}
