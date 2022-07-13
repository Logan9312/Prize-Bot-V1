package routers

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//go:embed embedded/*.html
var templateFS embed.FS

func BotStatus() {
	r := mux.NewRouter().StrictSlash(true)
	HandleRequests(r)
	log.Fatal(http.ListenAndServe(":8080", r))
}

type StatusOutput struct {
	Message string `json:"message"`
}

func HandleRequests(r *mux.Router) {
	r.HandleFunc("/success", Success)
	r.HandleFunc("/auction-bot/status", GetStatus).Methods("GET")
}

// GetStatus responds with the availability status of this service
func GetStatus(w http.ResponseWriter, r *http.Request) {
	status := StatusOutput{
		Message: "Bot is available",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	err := json.NewEncoder(w).Encode(status)
	if err != nil {
		fmt.Println("Error encoding: ", err.Error())
	}
}

func Success(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFS(templateFS, "embedded/*.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
