package plugin

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	"github.com/drone/drone-go/plugin/validator"
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
		description string
		router      *Router
		input       string
		want        string
	}{
		{
			description: "test without convert plugins",
			router:      NewRouter(""),
			input:       "name: default",
			want:        "name: default",
		},
		{
			description: "test with multiple convert plugins",
			router:      NewRouter("", WithConvertPlugins(newAddNewline(), newAddNewline())),
			input:       "name: default",
			want:        "name: default\n\n",
		},
	}

	for _, test := range tests {
		req := &converter.Request{Config: drone.Config{Data: test.input}}
		conf, err := test.router.Convert(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, test.want, conf.Data, test.description)
	}
}

// stringValidate returns an error if the passed in string is not in the validator requests Config.Data
type stringValidate struct {
	s string
}

func newStringValidate(s string) *stringValidate {
	return &stringValidate{s: s}
}

func (p *stringValidate) Validate(ctx context.Context, req *validator.Request) error {
	if !strings.Contains(req.Config.Data, p.s) {
		return errors.New(p.s)
	}

	return nil
}

func TestValidateRouter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		router      *Router
		input       string
		err         error
	}{
		{
			description: "test without validate plugins",
			router:      NewRouter(""),
		},
		{
			description: "test passing validate plugin",
			router:      NewRouter("", WithValidatePlugins(newStringValidate("one"))),
			input:       "one",
		},
		{
			description: "test with multiple validate plugins and an error",
			router:      NewRouter("", WithValidatePlugins(newStringValidate("one"), newStringValidate("two"))),
			input:       "one",
			err:         errors.New("two"),
		},
	}

	for _, test := range tests {
		req := &validator.Request{Config: drone.Config{Data: test.input}}
		assert.Equal(t, test.err, test.router.Validate(context.TODO(), req))
	}
}
