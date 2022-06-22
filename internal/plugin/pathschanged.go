package plugin

import (
	"context"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	pathschanged "github.com/meltwater/drone-convert-pathschanged/plugin"
)

type PathsChanged struct {
	provider string
}

func NewPathsChanged() *PathsChanged {
	return &PathsChanged{
		provider: "github",
	}
}

func (p *PathsChanged) Convert(ctx context.Context, req *converter.Request) (*drone.Config, error) {
	// TODO: this type can be removed once the upstream pathschanged
	// plugin supports reading the drone token directly
	plugin := pathschanged.New(p.provider, &pathschanged.Params{Token: req.Token.Access})

	return plugin.Convert(ctx, req)
}
