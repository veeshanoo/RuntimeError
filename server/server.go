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
	Address string
	Port    string
}

type ServerBuilder struct {
	config          *ServerConfig
	router          *mux.Router
	authService     services.AuthService
	questionService services.QuestionService
}

func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{}
}

func (b *ServerBuilder) WithConfig(config *ServerConfig) *ServerBuilder {
	b.config = config
	return b
}

func (b *ServerBuilder) WithRouter(router *mux.Router) *ServerBuilder {
	b.router = router
	return b
}

func (b *ServerBuilder) WithAuthService(authService services.AuthService) *ServerBuilder {
	b.authService = authService
	return b
}

func (b *ServerBuilder) WithQuestionService(questionService services.QuestionService) *ServerBuilder {
	b.questionService = questionService
	return b
}

func (b *ServerBuilder) Build() *Server {
	server := &Server{
		Config:          b.config,
		Router:          b.router,
		authService:     b.authService,
		questionService: b.questionService,
	}

	return server
}

type Server struct {
	Config          *ServerConfig
	Router          *mux.Router
	authService     services.AuthService
	questionService services.QuestionService
}

func NewServer() *Server {
	return NewServerBuilder().
		WithRouter(mux.NewRouter()).
		WithAuthService(services.NewAuthService()).
		WithQuestionService(services.NewQuestionService()).
		Build()
}

func (s *Server) Init() {
	s.BuildRoutes()
}

func (s *Server) Run() {
	var addr string
	if s.Config != nil {
		addr = fmt.Sprintf("%s:%s", s.Config.Address, s.Config.Port)
	} else {
		addr = "0.0.0.0:8080"
	}

	server := &http.Server{
		Addr:    addr,
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
	apiRouter.HandleFunc("/questions/create", s.CreateQuestion).Methods(http.MethodPost)
	apiRouter.HandleFunc("/questions/getAll", s.GetAll).Methods(http.MethodGet)
	apiRouter.HandleFunc("/questions/editContent", s.EditContent).Methods(http.MethodPut)

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
			"/api/login":    true,
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
