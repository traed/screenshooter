package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// Server - Handles API requests
type Server struct {
	httpServer *http.Server
	router     chi.Router
	savePath   string
}

// Start - Starts the Server
func (s *Server) Start() {
	s.savePath = "../uploads"
	s.router = chi.NewRouter()
	s.routes()
	s.httpServer = &http.Server{Addr: ":8080", Handler: s.router}

	log.Printf("Server ready on %s", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()

	if err != http.ErrServerClosed {
		log.Print("Http Server stopped unexpected")
		log.Print(err)
		s.Stop()
	} else {
		log.Print("Http Server stopped")
	}
}

// Stop - Attempts to gracefully stop the Server
func (s *Server) Stop() {
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			log.Print("Failed to shutdown http server gracefully")
		} else {
			s.httpServer = nil
		}
	}
}
