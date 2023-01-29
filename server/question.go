package server

import (
	types "RuntimeError/types/domain"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
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

	iface, _ := s.authService.Verify(context.TODO(), r.Header.Get("auth-token"))
	userId := iface.(string)
	question.SumitterId = userId

	questionId, err := s.questionService.CreateQuestion(context.Background(), question)
	if err != nil {
		respondWithError(w, err, http.StatusUnauthorized)
		return
	}

	respondWithJson(w, questionId)
}

func (s *Server) EditContent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	req := &types.EditContentRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		respondWithError(w, err, http.StatusBadRequest)
		return
	}

	iface, _ := s.authService.Verify(context.TODO(), r.Header.Get("auth-token"))
	userId := iface.(string)

	err = s.questionService.EditConent(context.Background(), userId, req)
	if err != nil {
		var status int
		switch err.Error() {
		case "Not found":
			status = http.StatusNotFound
		case "Unauthorized":
			status = http.StatusUnauthorized
		default:
			status = http.StatusInternalServerError
		}

		respondWithError(w, err, status)
		return
	}

	respondWithJson(w, nil)
}

func (s *Server) FavoriteAnswer(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	iface, _ := s.authService.Verify(context.TODO(), r.Header.Get("auth-token"))
	userId := iface.(string)

	favcomreq := &types.FavoriteCommentRequest{}
	if err := json.Unmarshal(body, favcomreq); err != nil {
		if err != nil {
			respondWithError(w, err, http.StatusBadRequest)
			return
		}
	}

	if err := s.questionService.FavoriteAnswer(context.Background(), userId, favcomreq.QuestionId, favcomreq.AnswerId); err != nil {
		respondWithError(w, err, http.StatusUnauthorized)
		return
	}

	respondWithJson(w, nil)
}

func (s *Server) UpdateQuestionVotes(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	iface, _ := s.authService.Verify(context.TODO(), r.Header.Get("auth-token"))
	userId := iface.(string)

	updateReq := &types.UpdateQuestionVotesRequest{}
	if err := json.Unmarshal(body, updateReq); err != nil {
		if err != nil {
			respondWithError(w, err, http.StatusBadRequest)
			return
		}
	}

	if updateReq.Type != "upvote" && updateReq.Type != "downvote" {
		respondWithError(w, errors.New("invalid type"), http.StatusBadRequest)
		return
	}

	switch updateReq.Type {
	case "upvote":
		err = s.questionService.UpvoteQuestion(context.Background(), updateReq.QuestionId, userId)
		if err != nil {
			respondWithError(w, err, http.StatusInternalServerError)
			return
		}
	case "downvote":
		err = s.questionService.DownvoteQuestion(context.Background(), updateReq.QuestionId, userId)
		if err != nil {
			respondWithError(w, err, http.StatusInternalServerError)
			return
		}
	}

	respondWithJson(w, nil)
}

func (s *Server) AddAnswerToQuestion(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	iface, _ := s.authService.Verify(context.TODO(), r.Header.Get("auth-token"))
	userId := iface.(string)

	req := &types.AddAnswerRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		respondWithError(w, err, http.StatusBadRequest)
		return
	}
	req.SubmitterId = userId

	err = s.questionService.AddAnswerToQuestion(context.Background(), req)
	if err != nil {
		respondWithError(w, err, http.StatusNotFound)
		return
	}

	respondWithJson(w, nil)
}

func (s *Server) AddReplyToAnswer(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	iface, _ := s.authService.Verify(context.TODO(), r.Header.Get("auth-token"))
	userId := iface.(string)

	req := &types.AddReplyRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		respondWithError(w, err, http.StatusBadRequest)
		return
	}
	req.SubmitterId = userId

	err = s.questionService.AddReplyToAnswer(context.Background(), req)
	if err != nil {
		respondWithError(w, err, http.StatusNotFound)
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

func (s *Server) GetQuestion(w http.ResponseWriter, r *http.Request) {
	questionId := mux.Vars(r)["questionId"]
	q, err := s.questionService.GetQuestion(context.Background(), questionId)
	if err != nil {
		respondWithError(w, err, http.StatusNotFound)
		return
	}
	respondWithJson(w, q)
}
