package rest

import (
	"os"
	"os/signal"
	"test-task-rest-api/internal/logger"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app    *fiber.App
	addr   string
	logger logger.Logger
}

func NewServer(app *fiber.App, addr string, logger logger.Logger) *Server {
	return &Server{
		app:    app,
		addr:   addr,
		logger: logger,
	}
}

func (s *Server) StartWithGracefulShutdown() {
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := s.app.Shutdown(); err != nil {
			s.logger.Errorf("Server is not shutting down! Reason: %v", err)
		}

		s.logger.Info("Server has successfully shut down!")

		close(idleConnsClosed)
	}()

	if err := s.app.Listen(s.addr); err != nil {
		s.logger.Fatalf("Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}
