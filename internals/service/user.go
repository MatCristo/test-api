package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MatCristo/test-api/internals/models"
)

var GetUser = func(id string) (*models.UserResponse, error) {

	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%s", id)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Erro ao consultar a API externa: %w", err)
	}
	defer resp.Body.Close()

	var user models.UserResponse
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("Erro ao decodificar resposta da API: %w", err)
	}

	return &user, err
}
