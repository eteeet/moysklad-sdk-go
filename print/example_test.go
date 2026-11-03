package print_test

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/print"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Печать отгрузки в pdf: 303 — файл готов и скачивается отдельным вызовом, 202 — SDK не ждёт.
func ExampleService_Export() {
	client, err := transport.NewClient(transport.Config{})
	if err != nil {
		panic(err)
	}
	printer, err := print.NewService(client)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // пример не обращается к API
	template := model.Reference{Meta: model.Meta{
		Href: "https://api.moysklad.ru/api/remap/1.2/entity/demand/metadata/embeddedtemplate/8a5bdbc2-ddb6-34f1-b359-05219426b979",
		Type: "embeddedtemplate", MediaType: "application/json"}}
	result, err := printer.Export(ctx, "demand", "a86708d2-f8d3-4e67-8f04-6101158da808",
		print.ExportRequest{Template: template, Extension: "pdf"})
	if err != nil {
		fmt.Println("отменено")
		return
	}
	if result.StatusCode == http.StatusSeeOther {
		if _, err := printer.Download(ctx, *result, os.Stdout); err != nil {
			fmt.Println(err)
		}
	}
	// Output: отменено
}
