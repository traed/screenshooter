package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"../pkg/capturer"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type message struct {
	URLs []string `json:"urls"`
}

func (s *Server) routes() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(30 * time.Second))

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hi!")
	})
	s.router.Route("/screenshot", func(r chi.Router) {
		r.Get("/{filename}", s.handleGetScreenshot())
		r.Get("/list", s.handleListScreenshot())
		r.Post("/", s.handleTakeScreenshot())
	})
}

func (s *Server) handleGetScreenshot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		filename := ctx.Value("filename").(string)
		fmt.Fprintf(w, "Get a screenshot: %s", filename)
	}
}

func (s *Server) handleListScreenshot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Get all screenshots")
	}
}

func (s *Server) handleTakeScreenshot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Taking screenshot")

		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var msg message
		err = json.Unmarshal(body, &msg)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		for _, u := range msg.URLs {
			url, err := url.ParseRequestURI(u)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			capturer := capturer.NewCapturer(s.savePath, url)
			go capturer.Execute()
		}

		fmt.Fprint(w, "Started capture")
	}
}
