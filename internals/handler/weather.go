package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/MatCristo/test-api/internals/service"
)

func WeatherHandler(client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cidade := r.URL.Query().Get("city")
		if cidade == "" {
			http.Error(w, "Cidade é necessário", http.StatusBadRequest)
			return
		}

		apiKey := os.Getenv("OPENWEATHER_API_KEY")
		weather, err := service.Weather(client, cidade, apiKey)
		if err != nil {
			http.Error(w, "Erro ao consultar o serviço de clima", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		mensagem := fmt.Sprintf("A cidade %s está com a temperatura de %.1f graus Celcius e o clima está com %s.",
			weather.Name, weather.Main.Temp, weather.Weather[0].Description)

		response := map[string]string{
			"mensagem": mensagem,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
