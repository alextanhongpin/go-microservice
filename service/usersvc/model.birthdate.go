package usersvc

import (
	"database/sql/driver"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const BirthDateLayout = "2006-01-02" // YYYY-MM-DD

type BirthDate struct {
	time.Time `json:",omitempty"`
}

func (b BirthDate) MarshalJSON() ([]byte, error) {
	if b.Time.IsZero() {
		return []byte(`""`), nil
	}
	return b.Time.MarshalJSON()
}

func (b *BirthDate) Scan(value interface{}) error {
	if value == nil {
		zap.L().Debug("value is nil")
		return nil
	}
	var v string
	switch t := value.(type) {
	case []byte:
		v = string(t)
	case string:
		v = string(t)
	}
	if v == "" {
		return nil
	}
	var err error
	b.Time, err = time.Parse(BirthDateLayout, v)
	if err != nil {
		return errors.Wrap(err, "scan BirthDate failed")
	}
	return nil
}

func (b BirthDate) Value() (driver.Value, error) {
	if b.Time.IsZero() {
		return "", nil
	}
	return b.Time.Format(BirthDateLayout), nil
}
