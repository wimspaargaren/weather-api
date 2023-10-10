// Package api provides some abstractions on top of echo to make it easier to use.
package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server is the API server.
type Server struct {
	e *echo.Echo
}

// NewServer creates a new API server.
func NewServer() *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	return &Server{
		e: e,
	}
}

// RegisterRoute registers a route.
func (s *Server) RegisterRoute(method, path string, handlerFunc echo.HandlerFunc) {
	s.e.Add(method, path, handlerFunc)
}

// Run runs the API server and handles graceful shutdown.
func (s *Server) Run() error {
	errChan := make(chan error, 1)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		// Start server
		err := s.e.Start(":8000")
		if errors.Is(http.ErrServerClosed, err) {
			log.Default().Println("server closed gracefully")
		} else if err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-c:
	}
	log.Default().Println("Shutting down server")
	err := s.e.Shutdown(context.Background())
	if err != nil {
		log.Default().Println("unable to shutdown server")
	}
	return nil
}
