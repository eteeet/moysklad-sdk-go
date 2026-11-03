package moysklad_test

import (
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	moysklad "github.com/eteeet/moysklad-sdk-go"
)

// Строки статуса, которые не совпадают с именем метода группы Entity.
var statusHandles = map[string]string{
	"characteristic":        "Entity.Variant",
	"cashier":               "Entity.RetailStore",
	"companysettings":       "Entity.Metadata",
	"states":                "Entity.Metadata",
	"entity/metadata":       "Entity.Metadata",
	"namedfilter":           "Entity.NamedFilter",
	"files":                 "Entity.Files",
	"images":                "Entity.Images",
	"context/employee":      "Context",
	"context/application":   "Context",
	"usersettings":          "Context",
	"companysubscription":   "Context",
	"turnover":              "Report.Turnover",
	"notification":          "Notification",
	"notification/settings": "Notification",
	"dashboard":             "Report.Dashboard",
	"byoperations":          "Report.ByOperations",
	"audit":                 "Audit",
	"async":                 "Async",
	"wizard":                "Wizard",
	"security/token":        "Security",
	"download":              "Transport",
}

// docs/status-entities.md: каждая отметка [x] ведёт к handle фасада, каждый handle отмечен,
// «Итого» совпадает с отметками разделов. Фиктивных или забытых отметок нет.
func TestStatusEntities(t *testing.T) {
	data, err := os.ReadFile("docs/status-entities.md")
	if err != nil {
		t.Fatal(err)
	}
	c, err := moysklad.NewClient(moysklad.Config{})
	if err != nil {
		t.Fatal(err)
	}
	marked := map[string]bool{}    // handle → есть отметка [x]
	byLower := map[string]string{} // "entity.product" → "Entity.Product"
	for _, group := range []map[string]reflect.Value{handles(c, "", "Entity", "Report"),
		handles(c.Entity(), "Entity."), handles(c.Report(), "Report.")} {
		for name := range group {
			marked[name] = false
			byLower[strings.ToLower(name)] = name
		}
	}

	type count struct{ done, total int }
	sections := map[string]*count{}
	var section *count
	row := regexp.MustCompile(`^- \[([ x])\] (\S+)`)
	for line := range strings.Lines(string(data)) {
		if title, ok := strings.CutPrefix(strings.TrimSpace(line), "## "); ok {
			section = &count{}
			sections[strings.ToLower(title)] = section
			continue
		}
		m := row.FindStringSubmatch(line)
		if m == nil || section == nil {
			continue
		}
		section.total++
		if m[1] != "x" {
			continue
		}
		section.done++
		handle, ok := statusHandles[m[2]]
		if !ok {
			handle = byLower["entity."+m[2]]
		}
		if _, exists := marked[handle]; !exists {
			t.Errorf("[x] %s has no facade handle", m[2])
			continue
		}
		marked[handle] = true
	}
	for name, done := range marked {
		if !done {
			t.Errorf("handle %s has no [x] in status-entities", name)
		}
	}

	total := regexp.MustCompile(`(справочники|документы|отчёты|прочее) (\d+)/(\d+)`)
	matches := total.FindAllStringSubmatch(string(data), -1)
	// Слово в «Итого» → заголовок раздела; отличается только «прочее».
	sectionOfTotal := map[string]string{"справочники": "справочники", "документы": "документы", "отчёты": "отчёты", "прочее": "прочие ресурсы"}
	if len(matches) != len(sectionOfTotal) {
		t.Fatalf("Итого: %v", matches)
	}
	for _, m := range matches {
		done, _ := strconv.Atoi(m[2])
		all, _ := strconv.Atoi(m[3])
		if got := sections[sectionOfTotal[m[1]]]; got == nil || got.done != done || got.total != all {
			t.Errorf("Итого %s %d/%d, rows %+v", m[1], done, all, got)
		}
	}
}
