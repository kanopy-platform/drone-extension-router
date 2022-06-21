package plugin

import (
	"context"
	"testing"

	"github.com/drone/drone-go/plugin/converter"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	t.Parallel()

	r := NewRouter(
		WithConvertPlugins(
			NewAddNewline(),
			NewAddNewline(),
		),
	)
	conf, err := r.Convert(context.Background(), &converter.Request{})
	assert.NoError(t, err)
	assert.Equal(t, "\n\n", conf.Data)
}
