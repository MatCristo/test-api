package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MatCristo/test-api/internals/handler"
)

func main() {

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	http.HandleFunc("/weather", handler.WeatherHandler(client))
	http.HandleFunc("/user", handler.GetUserHandler)

	fmt.Println("Servidor iniciado na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
