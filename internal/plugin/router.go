package plugin

import (
	"context"
	"net/http"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	"github.com/drone/drone-go/plugin/logger"
	"github.com/drone/drone-go/plugin/validator"
)

type Router struct {
	logger          logger.Logger
	convertPlugins  []converter.Plugin
	convertHandler  http.Handler
	validatePlugins []validator.Plugin
	validateHandler http.Handler
}

type RouterOption func(*Router)

func WithConvertPlugins(plugins ...converter.Plugin) RouterOption {
	return func(r *Router) {
		r.convertPlugins = append(r.convertPlugins, plugins...)
	}
}

func WithValidatePlugins(plugins ...validator.Plugin) RouterOption {
	return func(r *Router) {
		r.validatePlugins = append(r.validatePlugins, plugins...)
	}
}

func WithLogger(l logger.Logger) RouterOption {
	return func(r *Router) {
		r.logger = l
	}
}

func NewRouter(secret string, opts ...RouterOption) *Router {
	router := &Router{
		logger: logger.Discard(),
	}

	for _, opt := range opts {
		opt(router)
	}

	router.convertHandler = converter.Handler(router, secret, router.logger)
	router.validateHandler = validator.Handler(secret, router, router.logger)

	return router
}

func (r *Router) Convert(ctx context.Context, req *converter.Request) (*drone.Config, error) {
	for _, plugin := range r.convertPlugins {
		out, err := plugin.Convert(ctx, req)
		if err != nil {
			return nil, err
		}

		// modify the request object before it gets passed to the next plugin
		req.Config = *out
	}

	return &req.Config, nil
}

func (r *Router) Validate(ctx context.Context, req *validator.Request) error {
	for _, plugin := range r.validatePlugins {
		if err := plugin.Validate(ctx, req); err != nil {
			return err
		}
	}

	return nil
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.Header.Get("Accept") {
	case converter.V1:
		r.convertHandler.ServeHTTP(res, req)
	case validator.V1:
		r.validateHandler.ServeHTTP(res, req)
	default:
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}
