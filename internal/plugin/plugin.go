package plugin

import (
	"context"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
)

// TODO: addNewline is a mock plugin that should be removed once a real plugin in implemented
type addNewline struct{}

func NewAddNewline() converter.Plugin {
	return &addNewline{}
}

func (p *addNewline) Convert(ctx context.Context, req *converter.Request) (*drone.Config, error) {
	// TODO this should be modified or removed. For
	// demonstration purposes we show how you can ignore
	// certain configuration files by file extension.
	if strings.HasSuffix(req.Repo.Config, ".xyz") {
		// a nil response instructs the Drone server to
		// use the configuration file as-is, without
		// modification.
		return nil, nil
	}

	// get the configuration file from the request.
	config := req.Config.Data

	// TODO this should be modified or removed. For
	// demonstration purposes we make a simple modification
	// to the configuration file and add a newline.
	config = config + "\n"

	// returns the modified configuration file.
	return &drone.Config{
		Data: config,
	}, nil
}
