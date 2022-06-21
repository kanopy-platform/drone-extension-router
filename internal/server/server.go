package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/drone/drone-go/plugin/converter"
	"github.com/kanopy-platform/drone-convert/internal/plugin"
	"github.com/kanopy-platform/go-http-middleware/logging"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	router       *http.ServeMux
	secret       string
	pluginRouter *plugin.Router
}

// options
type ServerOptFunc func(*Server)

func WithPluginRouter(router *plugin.Router) func(*Server) {
	return func(s *Server) {
		s.pluginRouter = router
	}
}

func New(secret string, opts ...ServerOptFunc) http.Handler {
	s := &Server{
		router: http.NewServeMux(),
		secret: secret,
		pluginRouter: plugin.NewRouter(
			plugin.WithConvertPlugins(
				plugin.NewAddNewline(),
			),
		),
	}

	for _, opt := range opts {
		opt(s)
	}

	logrusMiddleware := logging.NewLogrus()

	s.router.HandleFunc("/", s.handlePlugins())
	s.router.HandleFunc("/healthz", s.handleHealthz())

	return logrusMiddleware.Middleware(s.router)
}

func (s *Server) handlePlugins() http.HandlerFunc {
	logger := log.StandardLogger()
	convertHandler := converter.Handler(s.pluginRouter, s.secret, logger)

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Header.Get("Accept") {
		case converter.V1:
			convertHandler.ServeHTTP(w, r)
		}
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
