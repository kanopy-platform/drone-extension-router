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

// options
type ServerOptFunc func(*Server)

func New(handler http.Handler, opts ...ServerOptFunc) http.Handler {
	s := &Server{
		router: http.NewServeMux(),
	}

	for _, opt := range opts {
		opt(s)
	}

	logrusMiddleware := logging.NewLogrus()

	s.router.Handle("/", handler)
	s.router.HandleFunc("/healthz", s.handleHealthz())

	return logrusMiddleware.Middleware(s.router)
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
