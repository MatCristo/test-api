package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestNewWeatherHandler_Sucesso(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("q") != "Serra" {
			t.Errorf("esperava cidade=Serra, obteve %s", query.Get("q"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{
            "name": "Serra",
            "main": { "temp": 29.5 },
            "weather": [ { "description": "ensolarado" } ]
        }`)
	}))
	defer mockServer.Close()

	client := &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			newReq, _ := http.NewRequest(req.Method, mockServer.URL, req.Body)
			newReq.URL.RawQuery = req.URL.RawQuery
			return mockServer.Client().Do(newReq)
		}),
	}

	os.Setenv("OPENWEATHER_API_KEY", "fake-key")

	req := httptest.NewRequest(http.MethodGet, "/weather?cidade=Serra", nil)
	rec := httptest.NewRecorder()

	handler := WeatherHandler(client)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperava status 200, veio %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Serra") || !strings.Contains(body, "29.5") {
		t.Errorf("resposta inesperada: %s", body)
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
