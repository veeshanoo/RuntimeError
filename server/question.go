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
	respondWithJson(w, question)

}

func (s *Server) FavoriteComment(w http.ResponseWriter, r *http.Request) {
	_, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	answerId := ""
	id := ""
	values := r.URL.Query()
	for k, v := range values {
		if k == "id" {
			id = v[0]
		} else if k == "answer" {
			answerId = v[0]
		}
	}
	_, err = s.questionService.FavoriteComment(context.Background(), id, answerId)
	if err != nil {
		respondWithError(w, err, http.StatusUnauthorized)
		return
	}
	respondWithJson(w, nil)

}

func (s *Server) DownvoteQuestion(w http.ResponseWriter, r *http.Request) {
	_, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	downvotterId := ""
	id := ""
	values := r.URL.Query()
	for k, v := range values {
		if k == "id" {
			id = v[0]
		} else if k == "downvotter" {
			downvotterId = v[0]
		}
	}
	_, err = s.questionService.DownvoteQuestion(context.Background(), id, downvotterId)
	if err != nil {
		respondWithError(w, err, http.StatusUnauthorized)
		return
	}
	respondWithJson(w, nil)

}

func (s *Server) UpvoteQuestion(w http.ResponseWriter, r *http.Request) {
	_, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	upvotterId := ""
	id := ""
	values := r.URL.Query()
	for k, v := range values {
		if k == "id" {
			id = v[0]
		} else if k == "upvotter" {
			upvotterId = v[0]
		}
	}
	_, err = s.questionService.UpvoteQuestion(context.Background(), id, upvotterId)
	if err != nil {
		respondWithError(w, err, http.StatusUnauthorized)
		return
	}
	respondWithJson(w, nil)

}

func (s *Server) AddAnswerToQuestion(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	values := r.URL.Query()
	id := ""
	for k, v := range values {
		if k == "id" {
			id = v[0]
		}
	}

	if err != nil {
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	answer := &types.Answer{}
	if err := json.Unmarshal(body, answer); err != nil {
		respondWithError(w, err, http.StatusBadRequest)
		return
	}
	_, err = s.questionService.AddAnswerToQuestion(context.Background(), id, answer)
	if err != nil {
		respondWithError(w, err, http.StatusUnauthorized)
		return
	}
	respondWithJson(w, id)

}

func (s *Server) AddReplyToAnswer(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	values := r.URL.Query()
	id := ""
	answerId := ""
	for k, v := range values {
		if k == "id" {
			id = v[0]
		} else if k == "answer" {
			answerId = v[0]
		}
	}

	if err != nil || id == "" || answerId == "" {
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	reply := &types.Reply{}
	if err := json.Unmarshal(body, reply); err != nil {
		respondWithError(w, err, http.StatusBadRequest)
		return
	}
	_, err = s.questionService.AddReplyToAnswer(context.Background(), id, answerId, reply)
	if err != nil {
		respondWithError(w, err, http.StatusUnauthorized)
		return
	}
	respondWithJson(w, id)

}

func (s *Server) GetAll(w http.ResponseWriter, r *http.Request) {

	questions, err := s.questionService.GetAll(context.Background())
	if err != nil {
		respondWithError(w, err, http.StatusUnauthorized)
		return
	}
	respondWithJson(w, questions)

}
