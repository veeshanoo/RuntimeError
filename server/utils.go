package server

import (
	"encoding/json"
	"net/http"
)

type errorStatus struct {
	Error string `json:"errorMessage"`
}

func respondWithError(w http.ResponseWriter, err error, statusCode int) {
	errStatus := errorStatus{
		Error: err.Error(),
	}

	w.WriteHeader(statusCode)
	if tmpErr := json.NewEncoder(w).Encode(errStatus); tmpErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}