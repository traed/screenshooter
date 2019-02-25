package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/traed/screenshooter/pkg/job"
)

type message struct {
	URLs []string `json:"urls"`
}

func (s *Server) routes() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	if s.ThrottleLimit > 0 {
		s.router.Use(middleware.Throttle(s.ThrottleLimit))
	}

	s.router.Use(middleware.Timeout(30 * time.Second))

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Go hack yourself!")
	})
	s.router.Route("/screenshot", func(r chi.Router) {
		r.Get("/{filename}", s.handleGetScreenshot())
		r.Post("/", s.handleTakeScreenshot())
	})
}

func (s *Server) handleGetScreenshot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := chi.URLParam(r, "filename")
		filepath := s.SavePath + "/" + filename
		pr, pw := io.Pipe()
		job := job.NewGetScreenshotJob(pw, filepath)

		s.Dispatcher.JobQueue <- &job

		content, err := ioutil.ReadAll(pr)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		if len(content) == 0 {
			http.Error(w, "File not found", 404)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(content)))
		if _, err := w.Write(content); err != nil {
			http.Error(w, "Failed writing response.", 500)
			return
		}
	}
}

func (s *Server) handleTakeScreenshot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var urls []string
		for _, u := range msg.URLs {
			url, err := url.ParseRequestURI(u)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			job := job.NewTakeScreenshotJob(s.SavePath, url)
			urls = append(urls, r.Host+r.URL.RequestURI()+"/"+job.GetFilename())

			s.Dispatcher.JobQueue <- &job
		}

		fmt.Fprintf(w, "%s", strings.Join(urls, "\n"))
	}
}
