package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type chirpData struct {
		Body string `json:"body"`
	}
	type cleanedResponse struct {
		Body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := chirpData{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode JSON", err)
		return
	}

	const maxChirpLength = 140
	if len(chirp.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {}, "sharbert": {}, "fornax": {},
	}
	cleanRes := cleanedResponse{Body: clean(chirp.Body, badWords)}

	respondWithJSON(w, http.StatusOK, cleanRes)
	return
}

func clean(message string, badWords map[string]struct{}) string {
	words := strings.Split(message, " ")

	for i, word := range words {
		if _, forbidden := badWords[strings.ToLower(word)]; forbidden {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
