package model

import (
	"database/sql/driver"
	"time"
)

const DateFormat = "2006-01-02"

type CustomDate time.Time

// used for dates in POST body
func (d *CustomDate) UnmarshalJSON(data []byte) error {
	// handle e
	parsed, err := time.Parse(`"`+DateFormat+`"`, string(data))
	if err != nil {
		return err
	}

	*d = CustomDate(parsed)
	return nil
}

// used for response body
func (d CustomDate) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(DateFormat)+2)
	b = append(b, '"')
	b = time.Time(d).AppendFormat(b, DateFormat)
	b = append(b, '"')
	return b, nil
}

// used when writing to db
func (d CustomDate) Value() (driver.Value, error) {
	if d.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}

	return []byte(time.Time(d).Format(DateFormat)), nil
}

// used when reading from db
func (d *CustomDate) Scan(v interface{}) error {
	// parse from postgres date format
	parsed, _ := time.Parse("2006-01-02 00:00:00 +0000 UTC", v.(time.Time).String())
	*d = CustomDate(parsed)
	return nil
}

func (d CustomDate) String() string {
	return time.Time(d).Format(DateFormat)
}

func (d CustomDate) Time() time.Time {
	return time.Time(d)
}
