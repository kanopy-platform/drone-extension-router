package plugin

import (
	"context"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	"github.com/stretchr/testify/assert"
)

type newline struct{}

func (p *newline) Convert(ctx context.Context, req *converter.Request) (*drone.Config, error) {
	return &drone.Config{Data: req.Config.Data + "\n"}, nil
}

func TestDefaultRouter(t *testing.T) {
	t.Parallel()

	r := NewRouter()

	req := &converter.Request{Config: drone.Config{Data: "name: default"}}
	assertConvert(t, r, "name: default", req)
}

func TestMultiRouter(t *testing.T) {
	t.Parallel()

	r := NewRouter(
		WithConvertPlugins(
			&newline{},
			&newline{},
		),
	)

	req := &converter.Request{Config: drone.Config{Data: "name: default"}}
	assertConvert(t, r, "name: default\n\n", req)
}

func assertConvert(t *testing.T, r *Router, want string, req *converter.Request) {
	conf, err := r.Convert(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, want, conf.Data)
}
