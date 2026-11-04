package vendorapi_test

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/vendorapi"
)

func fixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("../testdata/custom/vendorapi/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func decodeActivation(t *testing.T, codec documenttest.Codec, name string) vendorapi.ActivationRequest {
	t.Helper()
	var got vendorapi.ActivationRequest
	if err := codec.Unmarshal(fixture(t, name), &got); err != nil {
		t.Fatal(err)
	}
	return got
}

func TestActivationRequestFixtures(t *testing.T) {
	expiry := time.Date(2024, 1, 19, 18, 50, 12, 0, time.FixedZone("", 3*3600))
	subscription := func(trial bool, expiry time.Time) vendorapi.Subscription {
		return vendorapi.Subscription{TariffID: "23ca69d4-2657-40c4-8ba1-6ce24ddeac2e", Trial: trial, TariffName: "Basic", ExpiryMoment: &expiry}
	}
	checkSubscription := func(t *testing.T, got *vendorapi.Subscription, want vendorapi.Subscription) {
		t.Helper()
		if got == nil || got.TariffID != want.TariffID || got.Trial != want.Trial || got.TariffName != want.TariffName ||
			got.ExpiryMoment == nil || !got.ExpiryMoment.Equal(*want.ExpiryMoment) || got.NotForResale || got.Partner {
			t.Errorf("subscription = %+v", got)
		}
	}
	checkAdminAccess := func(t *testing.T, got []vendorapi.Access) {
		t.Helper()
		if len(got) != 1 || got[0].Resource != "https://api.moysklad.ru/api/remap/1.2" || len(got[0].Scope) != 1 || got[0].Scope[0] != "admin" ||
			got[0].AccessToken != "6ab89be1ae6ff147755625ee8da948e42612233b" || got[0].Permissions != nil {
			t.Errorf("access = %+v", got)
		}
	}
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			for name, cause := range map[string]vendorapi.Cause{
				"activation-install.json": vendorapi.CauseInstall, "activation-resume.json": vendorapi.CauseResume,
			} {
				got := decodeActivation(t, codec, name)
				if got.AppUID != "example-app.example-vendor" || got.AccountName != "dummyaccount" || got.Cause != cause || got.Additional != nil {
					t.Errorf("%s: %+v", name, got)
				}
				checkAdminAccess(t, got.Access)
				checkSubscription(t, got.Subscription, subscription(cause == vendorapi.CauseInstall, expiry))
			}

			got := decodeActivation(t, codec, "activation-permissions.json")
			if got.Cause != vendorapi.CauseInstall || got.AccountName != "account-test" || got.Subscription != nil || len(got.Access) != 1 ||
				got.Access[0].Scope[0] != "custom" || got.Access[0].AccessToken != "test-token" {
				t.Errorf("permissions: %+v", got)
			}
			documenttest.JSONEqual(t, got.Access[0].Permissions,
				[]byte(`{"supply":{"view":"ALL","update":"ALL"},"viewDashboard":true,"viewAudit":true}`))

			// При смене тарифа и автопродлении access не передаётся.
			for name, cause := range map[string]vendorapi.Cause{
				"activation-tariffchanged.json": vendorapi.CauseTariffChanged, "activation-autoprolongation.json": vendorapi.CauseAutoprolongation,
			} {
				got := decodeActivation(t, codec, name)
				if got.Cause != cause || got.Access != nil {
					t.Errorf("%s: %+v", name, got)
				}
				want := expiry
				if cause == vendorapi.CauseAutoprolongation {
					want = expiry.AddDate(0, 1, 0)
				}
				checkSubscription(t, got.Subscription, subscription(false, want))
			}

			got = decodeActivation(t, codec, "activation-fiscalapi.json")
			checkAdminAccess(t, got.Access)
			if got.Subscription != nil || got.Additional == nil || got.Additional.FiscalAPI == nil ||
				got.Additional.FiscalAPI.ID != "23ca69d4-2657-40c4-8ba1-6ce24ddeac2e" || got.Additional.FiscalAPI.Token != "asdfasdfgaerdfgqawefgqaergqa" {
				t.Errorf("fiscal: %+v", got)
			}
		})
	}
}

func TestDeactivationRequestFixtures(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for name, cause := range map[string]vendorapi.Cause{
			"deactivation-uninstall.json": vendorapi.CauseUninstall, "deactivation-suspend.json": vendorapi.CauseSuspend,
		} {
			var got vendorapi.DeactivationRequest
			if err := codec.Unmarshal(fixture(t, name), &got); err != nil {
				t.Fatal(err)
			}
			if got != (vendorapi.DeactivationRequest{Cause: cause, AppUID: "example-app.example-vendor", AccountName: "account-test"}) {
				t.Errorf("%s %s: %+v", codec.Name, name, got)
			}
		}
	}
}

// Ответ на активацию кодируется ровно в JSON примеров документации.
func TestActivationResponseFixtures(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for name, status := range map[string]vendorapi.Status{
			"response-settingsrequired.json": vendorapi.StatusSettingsRequired, "response-activated.json": vendorapi.StatusActivated,
		} {
			data, err := codec.Marshal(vendorapi.ActivationResponse{Status: status})
			if err != nil {
				t.Fatal(err)
			}
			documenttest.JSONEqual(t, data, fixture(t, name))
		}
	}
}

// Неизвестные значения cause и status не теряются: источник добавляет новые причины.
func TestOpenEnums(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		var got vendorapi.ActivationRequest
		if err := codec.Unmarshal([]byte(`{"cause":"PermissionsChanged"}`), &got); err != nil || got.Cause != vendorapi.CausePermissionsChanged {
			t.Errorf("%s: %+v, %v", codec.Name, got, err)
		}
		if err := codec.Unmarshal([]byte(`{"cause":"SomethingNew"}`), &got); err != nil || got.Cause != "SomethingNew" {
			t.Errorf("%s: %+v, %v", codec.Name, got, err)
		}
		data, err := codec.Marshal(vendorapi.ActivationResponse{Status: "Future"})
		if err != nil || string(data) != `{"status":"Future"}` {
			t.Errorf("%s: %s, %v", codec.Name, data, err)
		}
	}
	// Значения — дословно из таблиц статусов и причин документации.
	for got, want := range map[string]string{
		string(vendorapi.StatusActivating): "Activating", string(vendorapi.StatusActivationFailed): "ActivationFailed",
		string(vendorapi.StatusSettingsRequired): "SettingsRequired", string(vendorapi.StatusActivated): "Activated",
		string(vendorapi.StatusDeactivating): "Deactivating", string(vendorapi.StatusDeactivationFailed): "DeactivationFailed",
		string(vendorapi.StatusSuspended): "Suspended", string(vendorapi.CauseInstall): "Install", string(vendorapi.CauseResume): "Resume",
		string(vendorapi.CauseTariffChanged): "TariffChanged", string(vendorapi.CauseAutoprolongation): "Autoprolongation",
		string(vendorapi.CauseUninstall): "Uninstall", string(vendorapi.CauseSuspend): "Suspend",
		string(vendorapi.CausePermissionsChanged): "PermissionsChanged",
	} {
		if got != want {
			t.Errorf("constant %q, want %q", got, want)
		}
	}
}

// ActivationPattern направляет PUT (активация) и DELETE (деактивация) в обработчик ServeMux.
func TestActivationPattern(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("PUT "+vendorapi.ActivationPattern, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "activate "+r.PathValue("appId")+" "+r.PathValue("accountId"))
	})
	mux.HandleFunc("DELETE "+vendorapi.ActivationPattern, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "deactivate "+r.PathValue("appId")+" "+r.PathValue("accountId"))
	})
	const app, account = "5f3c5489-6a17-48b7-9fe5-b2000eb807fe", "f088b0a7-9490-4a57-b804-393163e7680f"
	for method, want := range map[string]string{http.MethodPut: "activate", http.MethodDelete: "deactivate"} {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest(method, "/api/moysklad/vendor/1.0/apps/"+app+"/"+account, nil))
		if rec.Code != http.StatusOK || rec.Body.String() != want+" "+app+" "+account {
			t.Errorf("%s: %d %q", method, rec.Code, rec.Body)
		}
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/api/moysklad/vendor/1.0/apps/"+app+"/"+account+"/status", nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("status path: %d", rec.Code)
	}
}

// Токены из тела активации не попадают в fmt, slog и JSON; значение достаётся явно string(s).
func TestSecretRedacted(t *testing.T) {
	got := decodeActivation(t, documenttest.Codecs[0], "activation-fiscalapi.json")
	const token, fiscal = "6ab89be1ae6ff147755625ee8da948e42612233b", "asdfasdfgaerdfgqawefgqaergqa"
	if string(got.Access[0].AccessToken) != token || string(got.Additional.FiscalAPI.Token) != fiscal {
		t.Fatalf("values = %q, %q", got.Access[0].AccessToken, got.Additional.FiscalAPI.Token)
	}
	var logs bytes.Buffer
	slog.New(slog.NewTextHandler(&logs, nil)).Info("text", "body", got)
	slog.New(slog.NewJSONHandler(&logs, nil)).Info("json", "body", got, "token", got.Access[0].AccessToken)
	out := []string{fmt.Sprint(got), fmt.Sprintf("%v %+v %#v %s %q", got, &got, got, got.Access[0].AccessToken, got.Additional.FiscalAPI.Token), logs.String()}
	// Ошибочный глагол fmt печатает значение в %!verb(...): Formatter закрывает и его.
	for _, verb := range []string{"%d", "%x", "%X", "%o", "%c", "%U", "%e", "%10.3s", "%#v"} {
		out = append(out, fmt.Sprintf(verb, got.Access[0].AccessToken), fmt.Sprintf(verb, *got.Additional.FiscalAPI))
	}
	for _, codec := range documenttest.Codecs {
		data, err := codec.Marshal(got)
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, string(data))
	}
	for _, s := range out {
		if strings.Contains(s, token) || strings.Contains(s, fiscal) || !strings.Contains(s, "redacted") {
			t.Errorf("leak or no marker: %s", s)
		}
	}
}
