package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MatCristo/test-api/internals/models"
)

func Weather(client *http.Client, cidade, apiKey string) (*models.WeatherResponse, error) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&lang=pt_br",
		cidade, apiKey,
	)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro na resposta da API: status %d", resp.StatusCode)
	}

	var weather models.WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON: %w", err)
	}

	return &weather, nil
}
