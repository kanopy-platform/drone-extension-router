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

func TestRouter(t *testing.T) {
	t.Parallel()

	r := NewRouter(
		WithConvertPlugins(
			&newline{},
			&newline{},
		),
	)
	conf, err := r.Convert(context.Background(), &converter.Request{})
	assert.NoError(t, err)
	assert.Equal(t, "\n\n", conf.Data)
}
