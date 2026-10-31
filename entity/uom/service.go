package uom

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
		return nil, fmt.Errorf("moysklad: uom service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/uom#3-poluchit-edinicu-izmereniya
func (s *Service) Get(ctx context.Context, id string) (*UOM, error) {
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: uom id must be a UUID")
	}
	return s.do(ctx, "get uom "+id, http.MethodGet, "entity/uom/"+id, nil)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any) (*UOM, error) {
	req, err := s.client.NewRequest(ctx, method, endpoint, nil, body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	var result UOM
	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if result.ID == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no id", operation)
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
