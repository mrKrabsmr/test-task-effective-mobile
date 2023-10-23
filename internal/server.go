package server

import (
	"embed"
	"github.com/go-chi/chi/v5"
	v1 "github.com/mrKrabsmr/commerce-edu-api/internal/api/v1"
	"github.com/mrKrabsmr/commerce-edu-api/internal/apps"
	"github.com/mrKrabsmr/commerce-edu-api/internal/configs"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"net/http"
)

//go:embed apps/persons/migrations
var EmbedMigrations embed.FS

type APIServer struct {
	config *configs.Config
	logger *logrus.Logger
	router *chi.Mux
}

func NewAPIServer() *APIServer {
	return &APIServer{
		config: core.GetConfig(),
		logger: core.GetLogger(),
		router: chi.NewRouter(),
	}
}

func (s *APIServer) configureRoutes() {
	s.router.Route("/api", func(r chi.Router) {
		v1.ConfigureRoutes(r)
	})
}

func (s *APIServer) StartMigrations() {
	conn := core.GetDB()

	goose.SetBaseFS(EmbedMigrations)

	if err := goose.SetDialect(s.config.DBDialect); err != nil {
		panic(err)
	}

	if err := goose.Up(conn.DB, "apps/persons/migrations"); err != nil {
		panic(err)
	}
}

func (s *APIServer) DownMigrations() {
	conn := core.GetDB()

	goose.SetBaseFS(EmbedMigrations)

	if err := goose.SetDialect(s.config.DBDialect); err != nil {
		panic(err)
	}

	if err := goose.Down(conn.DB, "apps/persons/migrations"); err != nil {
		panic(err)
	}
}

func (s *APIServer) Run() {
	s.configureRoutes()

	server := http.Server{
		Addr:    s.config.Address,
		Handler: s.router,
	}

	s.logger.Infof("STARTING SERVER AT %v", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
