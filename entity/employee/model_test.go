package employee_test

import (
	"os"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/entity/employee"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
)

func TestEmployeeFixture(t *testing.T) {
	data, err := os.ReadFile("../../testdata/upstream/employee.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got employee.Employee
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != "33333333-3333-3333-3333-333333333333" || got.UID != "admin@company" || got.LastName != "Иванов" ||
				got.ShortFio != "Иванов И. И." || got.FullName != "Иванов Иван Иванович" || got.INN != "123454567" ||
				got.Position != "Менеджер" || !got.Shared || got.Archived {
				t.Fatalf("employee = %+v", got)
			}
			if got.Salary == nil || got.Salary.Value != 1231231 {
				t.Fatalf("salary = %+v", got.Salary)
			}
			if len(got.Cashiers) != 1 || got.Cashiers[0].RetailStore == nil || got.Cashiers[0].Employee == nil ||
				got.Cashiers[0].ID != "634a96d3-0324-11f1-0a81-0727000001ea" {
				t.Fatalf("cashiers = %+v", got.Cashiers)
			}
			if got.Image == nil || got.Image.Size != 28593 || got.Image.Miniature.DownloadHref == "" || got.Image.Tiny.Href == "" ||
				time.Time(got.Image.Updated).IsZero() {
				t.Fatalf("image = %+v", got.Image)
			}
			if len(got.Attributes) != 3 || got.Owner == nil || got.Group == nil || time.Time(got.Created).IsZero() {
				t.Fatalf("attributes/owner/group/created = %+v", got)
			}
		})
	}
}

// Сотрудник без оклада, фото и кассиров: поля не выдумываются.
func TestEmployeeSparse(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		var got employee.Employee
		if err := codec.Unmarshal([]byte(`{"id":"x","lastName":"Л","unknownField":1}`), &got); err != nil {
			t.Fatal(err)
		}
		if got.Salary != nil || got.Image != nil || got.Cashiers != nil || got.Owner != nil {
			t.Fatalf("%s: invented fields = %+v", codec.Name, got)
		}
	}
}
