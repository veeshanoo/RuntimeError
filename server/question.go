package server

import (
	types "RuntimeError/types/domain"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (s *Server) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("haha0")
	body, err := io.ReadAll(r.Body)
	fmt.Println("haha1")
	if err != nil {
		fmt.Println("haha2")
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	question := &types.Question{}
	if err := json.Unmarshal(body, question); err != nil {
		fmt.Println("haha3")
		respondWithError(w, err, http.StatusBadRequest)
		return
	}
	fmt.Println("haha4")
	fmt.Println(question)
	_, err = s.questionService.CreateQuestion(context.Background(), question)
	if err != nil {
		fmt.Println("haha5")
		respondWithError(w, err, http.StatusUnauthorized)
		return
	}
	fmt.Println("haha6")
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
