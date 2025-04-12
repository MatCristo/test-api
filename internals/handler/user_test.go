package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MatCristo/test-api/internals/models"
	"github.com/MatCristo/test-api/internals/service"
)

func TestGetUserHandler(t *testing.T) {

	service.GetUser = func(id string) (*models.UserResponse, error) {
		return &models.UserResponse{
			Id:    1,
			Name:  "Breno",
			Email: "brenodev@dominioqualquer.com",
			Phone: "40028922",
			Address: struct {
				Street  string `json:"street"`
				Suite   string `json:"suite"`
				City    string `json:"city"`
				Zipcode string `json:"zipcode"`
			}{
				"Rua de baixo", "69", "Xique Xique", "70707-222",
			},
			Company: struct {
				Name        string `json:"name"`
				CatchPhrase string `json:"catchPhrase"`
				Bs          string `json:"bs"`
			}{
				"Empresa X", "", "",
			},
		}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/user?id=1", nil)
	rr := httptest.NewRecorder()

	GetUserHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status esperado %v, obtido %v", http.StatusOK, status)
	}

	var resp map[string]string
	json.NewDecoder(rr.Body).Decode(&resp)

	if resp["mensagem"] == "" {
		t.Errorf("esperava mensagem, mas veio vazia")
	}
}

func TestGetUserHandler_SemID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	rr := httptest.NewRecorder()

	GetUserHandler(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("esperado status %v, obtido %v", http.StatusBadRequest, status)
	}
}

func TestGetUserHandler_ErroService(t *testing.T) {
	service.GetUser = func(id string) (*models.UserResponse, error) {
		return nil, fmt.Errorf("simulando falha no serviço")
	}

	req := httptest.NewRequest(http.MethodGet, "/user?id=123", nil)
	rr := httptest.NewRecorder()

	GetUserHandler(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("esperado status %v, obtido %v", http.StatusInternalServerError, status)
	}
}
