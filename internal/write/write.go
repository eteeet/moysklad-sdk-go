// Package write — общий запрос удаления сервисов.
package write

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Delete — один DELETE endpoint. Успех — 200 или 204, тело ответа не читается. DELETE не
// повторяется: transport повторяет только GET и HEAD.
func Delete(ctx context.Context, client *transport.Client, operation, endpoint string, opts []transport.RequestOption) error {
	resp, err := client.Send(ctx, http.MethodDelete, endpoint, nil, nil, nil, opts...)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	return nil
}
