package server

import (
	"RuntimeError/services"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ServerConfig struct {
}

type Server struct {
	Config       *ServerConfig
	Router       *mux.Router
	loginService services.LoginService
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() {
	s.Router = mux.NewRouter()
	s.loginService = services.NewLoginService()
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
	apiRouter.Use(loggingMiddleware)
	apiRouter.Use(authMiddleware)

	// routes
	apiRouter.HandleFunc("/login", s.Login).Methods(http.MethodPost)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer next.ServeHTTP(w, r)

		// log request details
		fmt.Println(r.RequestURI)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer next.ServeHTTP(w, r)

		if r.URL.RequestURI() == "/api/login" {
			return
		}

		// check auth cookie
	})
}
