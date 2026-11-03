package notification_test

import (
	"context"
	"encoding/json/v2"
	"fmt"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/notification"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Поля конкретного типа уведомления читаются из Raw в свою структуру по meta.type.
func ExampleNotification_raw() {
	var n notification.Notification
	data := []byte(`{"meta":{"type":"NotificationImportCompleted"},"id":"0f423542-5b7a-11e9-9bea-3ff70000000f",` +
		`"read":true,"title":"Импорт завершен","taskType":"importer_good","taskState":"completed"}`)
	if err := json.Unmarshal(data, &n); err != nil {
		panic(err)
	}
	if n.Meta.Type == "NotificationImportCompleted" {
		var details struct {
			TaskType  string `json:"taskType"`
			TaskState string `json:"taskState"`
		}
		if err := json.Unmarshal(n.Raw, &details); err != nil {
			panic(err)
		}
		fmt.Println(n.Title, details.TaskType, details.TaskState)
	}
	// Output: Импорт завершен importer_good completed
}

// Отключить почту для задач и подписку на заказы покупателей; остальное сервер не меняет.
func ExampleService_UpdateSettings() {
	client, err := transport.NewClient(transport.Config{})
	if err != nil {
		panic(err)
	}
	service, err := notification.NewService(client)
	if err != nil {
		panic(err)
	}
	input := notification.SettingsRequest{
		"task":          {ChannelsEnabled: map[string]bool{"email": false}},
		"customerOrder": {Enabled: model.Value(false)},
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // пример не обращается к API
	if err := service.UpdateSettings(ctx, input); err != nil {
		fmt.Println("отменено")
	}
	// Output: отменено
}
