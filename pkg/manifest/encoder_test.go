package manifest_test

import (
	"strings"
	"testing"

	"github.com/kanopy-platform/drone-extension-router/pkg/manifest"
	"github.com/stretchr/testify/assert"
)

func TestEncoding(t *testing.T) {
	data := `
kind: pipeline
name: pipeline
---
kind: secret
name: secret
---
kind: signature
hmac: signature`

	resources, err := manifest.Decode(data)
	assert.NoError(t, err)
	assert.Equal(t, manifest.KindPipeline, resources[0].GetKind())
	assert.Equal(t, manifest.Kind("secret"), resources[1].GetKind())
	assert.Equal(t, manifest.Kind("signature"), resources[2].GetKind())

	pipeline, ok := resources[0].(*manifest.Pipeline)
	assert.True(t, ok)
	assert.Equal(t, "pipeline", pipeline.Name)

	encoded, err := manifest.Encode(resources)
	assert.NoError(t, err)

	docs := strings.Split(encoded, "\n---\n")
	assert.Equal(t, "---\nkind: pipeline\nname: pipeline\n", docs[0])
	assert.Equal(t, "kind: secret\nname: secret\n", docs[1])
	assert.Equal(t, "kind: signature\nhmac: signature\n", docs[2])
}
