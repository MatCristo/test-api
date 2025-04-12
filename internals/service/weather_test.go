package service

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeather_Sucesso(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		t.Logf("Mock server recebeu requisição: %s", r.URL.String())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{
            "name": "Serra",
            "main": { "temp": 27.5 },
            "weather": [ { "description": "céu limpo" } ]
        }`)
	}))
	defer mockServer.Close()

	client := &http.Client{
		Transport: &mockTransport{
			mockURL: mockServer.URL,
		},
	}

	testCidade := "Serra"
	testApiKey := "fakekey"

	weather, err := Weather(client, testCidade, testApiKey)

	if err != nil {
		t.Fatalf("esperava sucesso, mas deu erro: %v", err)
	}

	if weather == nil {
		t.Fatalf("esperava resposta, mas veio nil")
	}

	if weather.Name != "Serra" {
		t.Errorf("esperava cidade Serra, veio %s", weather.Name)
	}

	if weather.Main.Temp != 27.5 {
		t.Errorf("esperava temperatura 27.5, veio %.1f", weather.Main.Temp)
	}

	if len(weather.Weather) == 0 || weather.Weather[0].Description != "céu limpo" {
		t.Errorf("esperava descrição céu limpo, veio %+v", weather.Weather)
	}
}

type mockTransport struct {
	mockURL string
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	mockReq, _ := http.NewRequest(req.Method, m.mockURL, req.Body)

	return http.DefaultTransport.RoundTrip(mockReq)
}
