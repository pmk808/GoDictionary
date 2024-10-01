// main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Word struct {
	Text          string `json:"text"`
	Definition    string `json:"definition"`
	Pronunciation string `json:"pronunciation"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/word", getWord).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getWord(w http.ResponseWriter, r *http.Request) {
	word := r.URL.Query().Get("text")
	if word == "" {
		http.Error(w, "Missing 'text' parameter", http.StatusBadRequest)
		return
	}

	// TODO: Implement actual dictionary lookup
	result := Word{
		Text:          word,
		Definition:    "Sample definition",
		Pronunciation: "Sample pronunciation",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
