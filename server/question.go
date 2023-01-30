package server

import (
	"RuntimeError/auth"
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
	claims := iface.(*auth.JWTClaims)
	question.SubmitterId = claims.UserId
	question.SubmitterEmail = claims.Email

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
	claims := iface.(*auth.JWTClaims)

	err = s.questionService.EditContent(context.Background(), claims.UserId, req)
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
	claims := iface.(*auth.JWTClaims)

	favcomreq := &types.FavoriteCommentRequest{}
	if err := json.Unmarshal(body, favcomreq); err != nil {
		if err != nil {
			respondWithError(w, err, http.StatusBadRequest)
			return
		}
	}

	if err := s.questionService.FavoriteAnswer(context.Background(), favcomreq.QuestionId, favcomreq.AnswerId, claims); err != nil {
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
	claims := iface.(*auth.JWTClaims)

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
		err = s.questionService.UpvoteQuestion(context.Background(), updateReq.QuestionId, claims.UserId)
		if err != nil {
			respondWithError(w, err, http.StatusInternalServerError)
			return
		}
	case "downvote":
		err = s.questionService.DownvoteQuestion(context.Background(), updateReq.QuestionId, claims.UserId)
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
	claims := iface.(*auth.JWTClaims)

	req := &types.AddAnswerRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		respondWithError(w, err, http.StatusBadRequest)
		return
	}

	err = s.questionService.AddAnswerToQuestion(context.Background(), req, claims)
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
	claims := iface.(*auth.JWTClaims)

	req := &types.AddReplyRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		respondWithError(w, err, http.StatusBadRequest)
		return
	}

	err = s.questionService.AddReplyToAnswer(context.Background(), req, claims)
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

func (s *Server) GetQuestionsForUser(w http.ResponseWriter, r *http.Request) {
	iface, _ := s.authService.Verify(context.TODO(), r.Header.Get("auth-token"))
	claims := iface.(*auth.JWTClaims)

	questions, err := s.questionService.GetQuestionsForUser(context.Background(), claims.UserId)
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

func (s *Server) GetIncomingSuggestions(w http.ResponseWriter, r *http.Request) {
	iface, _ := s.authService.Verify(context.TODO(), r.Header.Get("auth-token"))
	claims := iface.(*auth.JWTClaims)

	sug, err := s.suggestionService.GetIncomingSuggestions(context.Background(), claims.UserId)
	if err != nil {
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	respondWithJson(w, sug)
}

func (s *Server) GetOutgoingSuggestions(w http.ResponseWriter, r *http.Request) {
	iface, _ := s.authService.Verify(context.TODO(), r.Header.Get("auth-token"))
	claims := iface.(*auth.JWTClaims)

	sug, err := s.suggestionService.GetOutgoingSuggestions(context.Background(), claims.UserId)
	if err != nil {
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	respondWithJson(w, sug)
}

func (s *Server) AddSugestion(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	iface, _ := s.authService.Verify(context.TODO(), r.Header.Get("auth-token"))
	claims := iface.(*auth.JWTClaims)

	su := &types.EditSuggestion{}
	if err := json.Unmarshal(body, su); err != nil {
		respondWithError(w, err, http.StatusBadRequest)
		return
	}

	su.SubmitterId = claims.UserId
	su.SubmitterEmail = claims.Email

	if id, err := s.suggestionService.CreateSuggestion(context.Background(), su); err != nil {
		respondWithError(w, err, http.StatusBadRequest)
		return
	} else {
		respondWithJson(w, id)
	}
}

func (s *Server) ApproveSuggestion(w http.ResponseWriter, r *http.Request) {
	suId := mux.Vars(r)["id"]

	iface, _ := s.authService.Verify(context.TODO(), r.Header.Get("auth-token"))
	claims := iface.(*auth.JWTClaims)

	if err := s.suggestionService.ApproveSuggestion(context.Background(), claims.UserId, suId); err != nil {
		respondWithError(w, err, http.StatusBadRequest)
		return
	}

	respondWithJson(w, nil)
}

func (s *Server) RejectSuggestion(w http.ResponseWriter, r *http.Request) {
	suId := mux.Vars(r)["id"]

	iface, _ := s.authService.Verify(context.TODO(), r.Header.Get("auth-token"))
	claims := iface.(*auth.JWTClaims)

	if err := s.suggestionService.RejectSuggestion(context.Background(), claims.UserId, suId); err != nil {
		respondWithError(w, err, http.StatusBadRequest)
		return
	}

	respondWithJson(w, nil)
}
