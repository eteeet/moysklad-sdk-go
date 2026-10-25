package model

import "encoding/json/v2"

type ResponseContext struct {
	Employee *Reference `json:"employee,omitzero"`
}

// employee приходит и как {"meta": {...}} (коллекции), и плоским meta (отчёты), пример /report/stock/all:
// https://github.com/moysklad/api-remap-1.2-doc/blob/b1c618b7daf63f599df1b7d11d2d656f1f82f8d0/md/reports/_report_stock.md#L209-L215
func (c *ResponseContext) UnmarshalJSON(data []byte) error {
	var value struct {
		Employee *struct {
			Meta
			WrappedMeta *Meta `json:"meta"`
		} `json:"employee"`
	}

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	c.Employee = nil
	if value.Employee == nil {
		return nil
	}

	meta := value.Employee.Meta
	if value.Employee.WrappedMeta != nil {
		meta = *value.Employee.WrappedMeta
	}

	c.Employee = &Reference{Meta: meta}
	return nil
}
