package model

// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-format-daty-i-vremeni

import (
	"encoding/json/v2"
	"fmt"
	"sync"
	"time"
)

const dateTimeLayout = "2006-01-02 15:04:05"

var moscowLocation = sync.OnceValue(func() *time.Location {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		// Без tzdata сохраняем работу SDK, но исторические смещения Москвы теряются.
		return time.FixedZone("MSK", 3*60*60)
	}
	return location
})

type DateTime time.Time

func (d DateTime) IsZero() bool {
	return time.Time(d).IsZero()
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	if time.Time(d).IsZero() {
		return []byte("null"), nil
	}
	value := time.Time(d).In(moscowLocation())
	if year := value.Year(); year < 0 || year > 9999 {
		return nil, fmt.Errorf("datetime year %d is outside 0000..9999", year)
	}
	return json.Marshal(value.Format(dateTimeLayout + ".000"))
}

func (d *DateTime) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	if value == "" {
		*d = DateTime{}
		return nil
	}
	var parsed time.Time
	var err error
	// RFC 3339 встречается в upstream-фикстурах; запись остаётся в формате API.
	if len(value) > 10 && value[10] == 'T' {
		parsed, err = time.Parse(time.RFC3339Nano, value)
	} else {
		parsed, err = time.ParseInLocation(dateTimeLayout, value, moscowLocation())
	}
	if err != nil {
		return err
	}
	*d = DateTime(parsed.In(moscowLocation()))
	return nil
}
