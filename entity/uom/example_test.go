package uom_test

import (
	"context"
	"fmt"

	"github.com/eteeet/moysklad-sdk-go/entity/uom"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Обычный вызов не меняется; опции и метаданные ответа передаются только при необходимости.
func ExampleService_Get() {
	client, _ := transport.NewClient(transport.Config{})
	service, _ := uom.NewService(client)
	ctx := context.Background()

	unit, err := service.Get(ctx, "77777777-7777-7777-7777-777777777777")
	_, _ = unit, err

	var info transport.ResponseInfo
	unit, err = service.Get(ctx, "77777777-7777-7777-7777-777777777777",
		transport.WithHeader("Accept-Language", "en"),
		transport.WithResponseInfo(&info))
	if err != nil {
		fmt.Println(info.StatusCode, info.Header.Get("X-Lognex-Auth"))
	}
}

// Variadic opts меняют тип method value: старую сигнатуру нужно обернуть.
func ExampleService_Get_methodValue() {
	client, _ := transport.NewClient(transport.Config{})
	service, _ := uom.NewService(client)

	var get func(context.Context, string, ...transport.RequestOption) (*uom.UOM, error) = service.Get
	legacy := func(ctx context.Context, id string) (*uom.UOM, error) { return service.Get(ctx, id) }
	_, _ = get, legacy
}
