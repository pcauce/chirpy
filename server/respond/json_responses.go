package respond

import (
	"encoding/json"
	"log"
	"net/http"
)

func WithError(w http.ResponseWriter, code int, msg string, err error) {
	log.Println(err)
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	WithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func WithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		WithError(w, http.StatusInternalServerError, "Error marshalling JSON", err)
		return
	}
}
