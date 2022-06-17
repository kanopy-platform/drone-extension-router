package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kanopy-platform/go-http-middleware/logging"
)

type Server struct {
	router *http.ServeMux
}

func New() http.Handler {
	s := &Server{router: http.NewServeMux()}
	logrusMiddleware := logging.NewLogrus()

	s.router.HandleFunc("/", s.handleRoot())
	s.router.HandleFunc("/healthz", s.handleHealthz())

	return logrusMiddleware.Middleware(s.router)
}

func (s *Server) handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello world")
	}
}

func (s *Server) handleHealthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := map[string]string{
			"status": "ok",
		}

		bytes, err := json.Marshal(status)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	}
}
