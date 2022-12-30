package server

import (
	types "RuntimeError/types/domain"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func (s *Server) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	question := &types.Question{}
	if err := json.Unmarshal(body, question); err != nil {
		respondWithError(w, err, http.StatusBadRequest)
		return
	}
	_, err = s.questionService.CreateQuestion(context.Background(), question)
	if err != nil {
		respondWithError(w, err, http.StatusUnauthorized)
		return
	}
	respondWithJson(w, nil)

}

func (s *Server) EditContent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	question := &types.Question{}
	if err := json.Unmarshal(body, question); err != nil {
		respondWithError(w, err, http.StatusBadRequest)
		return
	}
	_, err = s.questionService.EditConent(context.Background(), question.Id, question.Contents)
	if err != nil {
		respondWithError(w, err, http.StatusUnauthorized)
		return
	}
	respondWithJson(w, nil)

}

func (s *Server) GetAll(w http.ResponseWriter, r *http.Request) {

	questions, err := s.questionService.GetAll(context.Background())
	if err != nil {
		respondWithError(w, err, http.StatusUnauthorized)
		return
	}
	respondWithJson(w, questions)

}
