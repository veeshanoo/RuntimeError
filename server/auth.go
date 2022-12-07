package server

import (
	"RuntimeError/types"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	user := &types.UserLogin{}
	if err := json.Unmarshal(body, user); err != nil {
		respondWithError(w, err, http.StatusBadRequest)
		return
	}

	token, err := s.authService.Login(context.Background(), user)
	if err != nil {
		respondWithError(w, err, http.StatusUnauthorized)
		return
	}

	respondWithJson(w, types.UserLoginResponse{Token: token})
}

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	user := &types.UserLogin{}
	if err := json.Unmarshal(body, user); err != nil {
		respondWithError(w, err, http.StatusBadRequest)
		return
	}

	fmt.Printf("%+v\n", user)

	err = s.authService.Register(context.Background(), user)
	if err != nil {
		respondWithError(w, err, http.StatusUnauthorized)
		return
	}

	respondWithJson(w, nil)
}
