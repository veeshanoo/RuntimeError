package server

import (
	"RuntimeError/services"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ServerConfig struct {
}

type Server struct {
	Config      *ServerConfig
	Router      *mux.Router
	authService services.AuthService
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() {
	s.Router = mux.NewRouter()
	s.authService = services.NewAuthService()
	s.BuildRoutes()
}

func (s *Server) Run() {
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: s.Router,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func (s *Server) BuildRoutes() {
	apiRouter := s.Router.PathPrefix("/api").Subrouter()

	// add logging and auth middleware
	apiRouter.Use(s.loggingMiddleware, s.authMiddleware)

	// routes
	apiRouter.HandleFunc("/login", s.Login).Methods(http.MethodPost)
	apiRouter.HandleFunc("/register", s.Register).Methods(http.MethodPost)
	apiRouter.HandleFunc("/test", s.Test).Methods(http.MethodGet)
}

func (s *Server) Test(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HIT TEST ENDPOINT")
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer next.ServeHTTP(w, r)

		// log request details
		fmt.Println(r.RequestURI)
	})
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		excludePaths := map[string]bool{
			"/api/login": true,
			"/api/register": true,
		}
		
		if _, ok := excludePaths[r.URL.RequestURI()]; ok {
			next.ServeHTTP(w, r)
			return
		}

		// check auth cookie
		token := r.Header.Get("auth-token")
		if token == "" {
			respondWithError(w, errors.New("invalid cookie"), http.StatusUnauthorized)
			return
		}

		if _, err := s.authService.Verify(context.TODO(), token); err != nil {
			respondWithError(w, errors.New("invalid cookie"), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
