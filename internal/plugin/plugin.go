package plugin

import (
	"context"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
)

type Router struct {
	convertPlugins []converter.Plugin
}

type RouterOption func(*Router)

func WithConvertPlugins(plugins ...converter.Plugin) RouterOption {
	return func(r *Router) {
		r.convertPlugins = append(r.convertPlugins, plugins...)
	}
}

func NewRouter(opts ...RouterOption) *Router {
	router := &Router{}

	for _, opt := range opts {
		opt(router)
	}

	return router
}

func (r *Router) Convert(ctx context.Context, req *converter.Request) (conf *drone.Config, err error) {
	for _, plugin := range r.convertPlugins {
		conf, err = plugin.Convert(ctx, req)
		if err != nil {
			return nil, err
		}

		// modify the request object before it gets passed to the next plugin
		req.Config = *conf
	}

	return
}
