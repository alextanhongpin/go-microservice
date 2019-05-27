package user

import "time"

type BirthDate struct {
	time.Time `json:"birth_date,omitempty"`
}

func (d BirthDate) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return d.Time.MarshalJSON()
}

func (d BirthDate) Value() time.Time {
	return d.Time
}
