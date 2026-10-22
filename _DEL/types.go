package moysklad

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// moscow — часовой пояс дат МойСклад. С 2014 года это постоянно UTC+3, поэтому tzdata не нужна.
func moscow() *time.Location { return time.FixedZone("MSK", 3*60*60) }

// timeLayout — формат дат МойСклад в ответах и фильтрах.
const timeLayout = "2006-01-02 15:04:05.000"

// Timestamp — дата МойСклад в формате "yyyy-MM-dd HH:mm:ss.SSS" по московскому времени.
type Timestamp struct{ time.Time }

// UnmarshalJSON разбирает дату с миллисекундами или без них. Принимает и ISO 8601 со смещением:
// так записаны даты в официальных фикстурах МойСклад (testdata/README.md).
func (ts *Timestamp) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	if strings.Contains(s, "T") {
		v, err := time.Parse(time.RFC3339Nano, s)
		if err != nil {
			return fmt.Errorf("time.Parse: %w", err)
		}

		ts.Time = v

		return nil
	}

	v, err := time.ParseInLocation(timeLayout[:min(len(s), len(timeLayout))], s, moscow())
	if err != nil {
		return fmt.Errorf("time.ParseInLocation: %w", err)
	}

	ts.Time = v

	return nil
}

// Attribute — дополнительное поле. Тип Value зависит от Type.
type Attribute struct {
	Meta  *Meta           `json:"meta,omitempty"`
	ID    string          `json:"id,omitempty"`
	Name  string          `json:"name,omitempty"`
	Type  string          `json:"type,omitempty"`
	Value json.RawMessage `json:"value,omitempty"`
}
