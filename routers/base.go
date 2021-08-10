package routers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)
func BotStatus () {
r := mux.NewRouter().StrictSlash(true)
HandleRequests(r)
log.Fatal(http.ListenAndServe(":8080", r))
}

type StatusOutput struct {
	Message string `json:"message"`
}

func HandleRequests(r *mux.Router) {
	r.HandleFunc("/auction-bot/status", GetStatus).Methods("GET")
}

// GetStatus responds with the availability status of this service
func GetStatus(w http.ResponseWriter, r *http.Request) {
	status := StatusOutput{
		Message: "Bot is available",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(status)
}