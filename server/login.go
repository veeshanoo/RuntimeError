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
		respondWithError(w, err, http.StatusInternalServerError)
	}

	fmt.Printf("%+v\n", user)

	usr, err := s.loginService.Register(context.TODO())
	if err != nil {
		respondWithError(w, err, http.StatusInternalServerError)
	}
	fmt.Printf("%+v\n", usr)
}