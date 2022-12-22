package plugin

import (
	"context"
	"log"

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
	log.Printf("YUZHOU DEBUG before conversion")
	log.Printf(req.Config.Data)

	// TODO: this type can be removed once the upstream pathschanged
	// plugin supports reading the drone token directly
	plugin := pathschanged.New(p.provider, &pathschanged.Params{Token: req.Token.Access})

	config, err := plugin.Convert(ctx, req)

	log.Printf("YUZHOU DEBUG after conversion")
	log.Printf(req.Config.Data)

	return config, err
}
