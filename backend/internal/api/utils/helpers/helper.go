package helpers

import (
	"net/http"
	"encoding/json"
	"log"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set(
		"Content-Type",
		"encoding/json",
	)

	w.WriteHeader(code)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed encoding with error %v\n", err)
	}
}

func WriteError(w http.ResponseWriter, code int, msg string) {
	resp := ErrorResponse{
		Error : msg,
	}
	WriteJSON(w, code, resp)
}	