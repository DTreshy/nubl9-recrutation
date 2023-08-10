package server

import (
	"github.com/DTreshy/nubl9-recrutation/api"
	json "github.com/goccy/go-json"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

type Server struct {
	app    *fiber.App
	logger *zap.Logger
}

func New(logger *zap.Logger) Server {
	server := Server{
		app: fiber.New(fiber.Config{
			JSONEncoder: json.Marshal,
			JSONDecoder: json.Unmarshal,
		}),
		logger: logger,
	}

	return server
}

func (s *Server) Run(bind string) error {
	s.logger.Sugar().Info("Starting HTTP server")
	return s.app.Listen(bind)
}

func (server *Server) Routes() error {
	server.app.Use(cors.New())

	server.app.Get("/random/+", api.Random)

	return nil
}

func (server *Server) Shutdown() error {
	server.logger.Sugar().Info("Stopping HTTP server")
	return server.app.Shutdown()
}
