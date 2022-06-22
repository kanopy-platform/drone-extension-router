package plugin

import (
	"context"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	"github.com/stretchr/testify/assert"
)

type addNewline struct{}

func newAddNewline() *addNewline {
	return &addNewline{}
}

func (p *addNewline) Convert(ctx context.Context, req *converter.Request) (*drone.Config, error) {
	return &drone.Config{Data: req.Config.Data + "\n"}, nil
}

func TestConvertRouter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		router *Router
		input  string
		want   string
	}{
		// test without convert plugins
		{
			router: NewRouter(),
			input:  "name: default",
			want:   "name: default",
		},
		// test with multiple convert plugins
		{
			router: NewRouter(WithConvertPlugins(newAddNewline(), newAddNewline())),
			input:  "name: default",
			want:   "name: default\n\n",
		},
	}

	for _, test := range tests {
		req := &converter.Request{Config: drone.Config{Data: test.input}}
		conf, err := test.router.Convert(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, test.want, conf.Data)
	}
}
