package errors_test

import (
	"crypto/sha256"
	"encoding/json/v2"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Каталог сгенерирован из _errors.md revision 9a811cef: 545 уникальных кодов, имя константы
// совпадает с числом, семантических алиасов нет. Дубли сообщений в источнике: 1004 и 1083
// («Ошибка при формировании результата») — разные коды, обе константы сохранены.
func TestCodesInventory(t *testing.T) {
	data, err := os.ReadFile("codes.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_errors.md",
		"709ead9554e75cfaf6430dcfdb0786ddb95bf0b17a1fe48c6c4a10ae7aff9199"} {
		if !strings.Contains(string(data), want) {
			t.Fatalf("provenance %s missing", want)
		}
	}
	file, err := parser.ParseFile(token.NewFileSet(), "codes.go", data, 0)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[int]bool{}
	var codes []string
	for _, decl := range file.Decls {
		for _, spec := range decl.(*ast.GenDecl).Specs {
			value := spec.(*ast.ValueSpec)
			name := value.Names[0].Name
			code, err := strconv.Atoi(value.Values[0].(*ast.BasicLit).Value)
			if err != nil || name != fmt.Sprintf("Code%d", code) || seen[code] {
				t.Fatalf("%s = %v", name, value.Values[0])
			}
			seen[code] = true
			codes = append(codes, strconv.Itoa(code))
		}
	}
	// SHA-256 отсортированного списка id="error_N" источника, посчитанный при генерации.
	sum := fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(codes, ","))))
	if len(codes) != 545 || sum != "98f26f5a0eaaf40621f3f8a134d918d01a435c10c765dd529379851891ddf4da" {
		t.Fatalf("codes = %d, sha256 = %s", len(codes), sum)
	}

	// Константы нетипизированные: присваиваются Detail.Code int без преобразования.
	detail := sdkerrors.Detail{Code: sdkerrors.Code1073}
	if detail.Code != 1073 || sdkerrors.Code1062 != 1062 || sdkerrors.Code3042 != 3042 {
		t.Fatal("constant values changed")
	}
}

func TestHasCodeAndStatus(t *testing.T) {
	var apiErr sdkerrors.APIError
	// Неизвестный будущий код и detail без code сохраняются при декодировании.
	if err := json.Unmarshal([]byte(`{"errors":[{"error":"new","code":99999},{"error":"no code"},{"error":"slow","code":1073}]}`), &apiErr); err != nil {
		t.Fatal(err)
	}
	apiErr.StatusCode = http.StatusTooManyRequests

	wrapped := fmt.Errorf("get uom: %w", &apiErr)
	joined := errors.Join(errors.New("other"), &sdkerrors.APIError{StatusCode: 429, Errors: []sdkerrors.Detail{{Code: 1000}}}, wrapped)
	for name, err := range map[string]error{"direct": &apiErr, "wrapped": wrapped, "joined": joined} {
		if !sdkerrors.HasCode(err, 99999) || !sdkerrors.HasCode(err, sdkerrors.Code1073) || sdkerrors.HasCode(err, sdkerrors.Code1001) ||
			sdkerrors.HasCode(err, sdkerrors.Code1000) != (name == "joined") {
			t.Errorf("%s: HasCode wrong", name)
		}
		if sdkerrors.StatusCode(err) != http.StatusTooManyRequests {
			t.Errorf("%s: StatusCode = %d", name, sdkerrors.StatusCode(err))
		}
	}
	// Код 0 не совпадает с detail без code.
	if apiErr.HasCode(0) || sdkerrors.HasCode(errors.New("x"), 0) || sdkerrors.HasCode(nil, 1000) {
		t.Error("HasCode on missing code")
	}
	if sdkerrors.StatusCode(nil) != 0 || sdkerrors.StatusCode(errors.New("x")) != 0 {
		t.Error("StatusCode without APIError")
	}
}

// 413 возвращается как APIError с деталями и заголовками за один запрос, без разбиения и повтора.
func TestPayloadTooLarge(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Header().Set("X-Lognex-Content-Region", "ru")
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		w.Write([]byte(`{"errors":[{"error":"too many","code":2007}]}`))
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	var info transport.ResponseInfo
	_, err = client.Send(t.Context(), http.MethodPost, "entity/demand", nil, map[string]any{"positions": []int{1}}, nil,
		transport.WithResponseInfo(&info))
	if !sdkerrors.HasCode(err, sdkerrors.Code2007) || sdkerrors.StatusCode(err) != http.StatusRequestEntityTooLarge ||
		strings.Count(err.Error(), "moysklad:") != 1 {
		t.Fatalf("err = %v", err)
	}
	if calls.Load() != 1 || info.Header.Get("X-Lognex-Content-Region") != "ru" {
		t.Fatalf("calls = %d, info = %+v", calls.Load(), info)
	}
}
