package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gradlink/internal/db"
	apphttp "gradlink/internal/http"
)

type TestAPI struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func NewAPI(t *testing.T) *TestAPI {
	t.Helper()

	gin.SetMode(gin.TestMode)

	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	router := gin.New()
	apphttp.RegisterRoutes(router, &apphttp.Deps{DB: database})

	return &TestAPI{Router: router, DB: database}
}

func (api *TestAPI) JSON(t *testing.T, method, path string, body any, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	t.Helper()

	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		payload, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
		reader = bytes.NewReader(payload)
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	rec := httptest.NewRecorder()
	api.Router.ServeHTTP(rec, req)
	return rec
}

func DecodeJSON(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode JSON response %q: %v", rec.Body.String(), err)
	}

	return payload
}

func ErrorCode(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()

	payload := DecodeJSON(t, rec)
	errPayload, ok := payload["error"].(map[string]any)
	if !ok {
		t.Fatalf("missing error object in response: %s", rec.Body.String())
	}
	code, ok := errPayload["code"].(string)
	if !ok {
		t.Fatalf("missing error code in response: %s", rec.Body.String())
	}

	return code
}

func SessionCookie(t *testing.T, rec *httptest.ResponseRecorder) *http.Cookie {
	t.Helper()

	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "gradlink_session" {
			return cookie
		}
	}

	t.Fatalf("missing session cookie in response: %s", rec.Body.String())
	return nil
}
