package http

import (
	"fmt"
	"lamoda-tech-assigment/internal/config"
	v1 "lamoda-tech-assigment/internal/entrypoints/http/v1"
	"lamoda-tech-assigment/internal/services/usecases"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	httpSwag "github.com/swaggo/http-swagger/v2"
)

type Server struct {
	srv *http.Server
	uc  usecases.IUseCases
}

// NewServer returns pointer to configured http server
func NewServer(cfg config.HTTP, uc usecases.IUseCases) *Server {

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      nil,
	}

	return &Server{srv: srv, uc: uc}
}

func (s *Server) ConfigureRouter() {
	router := mux.NewRouter().StrictSlash(true)

	api := router.PathPrefix("/api").Subrouter()

	handlerV1 := v1.NewHandler(s.uc)
	handlerV1.InitializeRoutes(api.PathPrefix("/v1").Subrouter())

	router.PathPrefix("/swagger").Handler(httpSwag.Handler(
		httpSwag.URL("http://localhost:8080/swagger/doc.json"),
		httpSwag.DeepLinking(true),
		httpSwag.DocExpansion("list"),
		httpSwag.DomID("swagger-ui"),
	))

	s.srv.Handler = router
}

func (s *Server) ListenAndServe() error {
	if s.srv.Handler == nil {
		s.ConfigureRouter()
	}

	return s.srv.ListenAndServe()
}
