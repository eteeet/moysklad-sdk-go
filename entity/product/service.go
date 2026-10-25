package product

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/eteeet/moysklad-sdk-go/transport"
)

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: product service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#3-poluchit-tovar
func (s *Service) Get(ctx context.Context, id string) (*Product, error) {
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: product id must be a UUID")
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, "entity/product/"+id, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("moysklad: get product: %w", err)
	}

	var result Product
	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, fmt.Errorf("moysklad: get product: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: get product: unexpected HTTP status %d", resp.StatusCode)
	}
	if result.ID == "" {
		return nil, fmt.Errorf("moysklad: get product: response has no id")
	}
	return &result, nil
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}
