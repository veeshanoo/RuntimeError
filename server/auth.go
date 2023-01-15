package server

import (
	"RuntimeError/types/domain"
	"context"
	"encoding/json"
	"errors"
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

	user := &types.User{}
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

	user := &types.User{}
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

func (s *Server) GetUserData(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("auth-token")
	idIface, err := s.authService.Verify(context.TODO(), token)
	if err != nil {
		respondWithError(w, errors.New("invalid cookie"), http.StatusUnauthorized)
		return
	}

	userId := idIface.(string)
	userData, err := s.authService.GetUserData(context.TODO(), userId)
	if err != nil {
		respondWithError(w, errors.New("unknown"), http.StatusInternalServerError)
		return
	}

	respondWithJson(w, userData)
}
