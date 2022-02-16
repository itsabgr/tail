package handy

import (
	"bytes"
	"time"
)

//ZeroTime returns a zero time
func ZeroTime() *time.Time {
	return &time.Time{}
}

//IsZeroTime checks aTime is zero
func IsZeroTime(aTime *time.Time) bool {
	return time.Time{} == *aTime
}

//Time is a nullable time
type Time struct {
	time.Time
	//IsFill is true when is not nil
	IsSet bool `json:"is_set"`
}

func (time *Time) MarshalJSON() ([]byte, error) {
	if time.IsSet {
		return time.MarshalJSON()
	}
	return []byte("null"), nil
}

func (time *Time) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		time.IsSet = false
		time.Time = *ZeroTime()
		return nil
	}
	return time.UnmarshalJSON(b)
}
