package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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

		_, err := os.Stat(filepath)
		if os.IsNotExist(err) {
			http.Error(w, "File not found.", 404)
			return
		} else if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		file, err := os.Open(filepath)
		defer file.Close()
		if err != nil {
			http.Error(w, "Unable to open file.", 500)
			return
		}

		img, err := png.Decode(file)
		if err != nil {
			http.Error(w, "Unable to decode file.", 500)
			return
		}

		buffer := new(bytes.Buffer)
		err = png.Encode(buffer, img)
		if err != nil {
			http.Error(w, "Unable to create file buffer.", 500)
			return
		}

		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		if _, err := w.Write(buffer.Bytes()); err != nil {
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

			capturer := capturer.NewCapturer(s.SavePath, url)
			urls = append(urls, r.Host+r.URL.RequestURI()+"/"+capturer.GetFilename())

			go capturer.Execute()
		}

		fmt.Fprintf(w, "%s", strings.Join(urls, "\n"))
	}
}
