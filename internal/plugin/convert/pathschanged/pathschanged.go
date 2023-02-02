package pathschanged

import (
	"context"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	"github.com/meltwater/drone-convert-pathschanged/plugin"
)

type PathsChanged struct {
	provider string
}

func New() *PathsChanged {
	return &PathsChanged{
		provider: "github",
	}
}

func (p *PathsChanged) Convert(ctx context.Context, req *converter.Request) (*drone.Config, error) {
	// TODO: this type can be removed once the upstream pathschanged
	// plugin supports reading the drone token directly
	plugin := plugin.New(p.provider, &plugin.Params{Token: req.Token.Access})

	return plugin.Convert(ctx, req)
}
