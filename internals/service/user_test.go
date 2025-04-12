package service

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var originalAPI = "https://jsonplaceholder.typicode.com"

func TestGetUser_Sucesso(t *testing.T) {
	mockResponse := `{
		"id": 1,
		"name": "Breno",
		"email": "brenodev@dominioqualquer.com",
		"phone": "40028922",
		"address": {
			"street": "Rua de baixo",
			"suite": "69",
			"city": "Xique Xique",
			"zipcode": "70707-222"
		},
		"company": {
			"name": "Empresa X",
			"catchPhrase": "",
			"bs": ""
		}
	}`

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, mockResponse)
	}))
	defer mockServer.Close()

	originalTransport := http.DefaultTransport
	defer func() { http.DefaultTransport = originalTransport }()

	http.DefaultTransport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if strings.Contains(req.URL.Host, "jsonplaceholder.typicode.com") {
			req.URL.Scheme = "http"
			req.URL.Host = strings.TrimPrefix(mockServer.URL, "http://")
		}
		return originalTransport.RoundTrip(req)
	})

	user, err := GetUser("1")
	if err != nil {
		t.Fatalf("esperava sucesso, mas deu erro: %v", err)
	}
	if user.Name != "Breno" {
		t.Errorf("esperava nome Breno, veio %s", user.Name)
	}
}

func TestGetUser_ErroAPIRetorna500(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Erro interno", http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	originalTransport := http.DefaultTransport
	defer func() { http.DefaultTransport = originalTransport }()

	http.DefaultTransport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		req.URL.Scheme = "http"
		req.URL.Host = strings.TrimPrefix(mockServer.URL, "http://")
		return originalTransport.RoundTrip(req)
	})

	_, err := GetUser("1")
	if err == nil {
		t.Fatal("esperava erro da API, mas veio nil")
	}
}
