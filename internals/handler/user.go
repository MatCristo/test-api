package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MatCristo/test-api/internals/service"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "ID de usuário é necessário", http.StatusBadRequest)
		return
	}

	user, err := service.GetUser(id)
	if err != nil {
		http.Error(w, "Erro ao consultar o serviço de usuarios", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	mensagem := fmt.Sprintf(
		"O usuário com id %d é %s (Email:%s, Telefone %s) vive na rua %s numero %s na cidade de %s com o código postal de %s e trabalha na empresa %s.",
		user.Id,
		user.Name,
		user.Email,
		user.Phone,
		user.Address.Street,
		user.Address.Suite,
		user.Address.City,
		user.Address.Zipcode,
		user.Company.Name,
	)

	response := map[string]string{
		"mensagem": mensagem,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
