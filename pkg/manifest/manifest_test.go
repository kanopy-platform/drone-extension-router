package manifest_test

import (
	"strings"
	"testing"

	"github.com/kanopy-platform/drone-extension-router/pkg/manifest"
	"github.com/stretchr/testify/assert"
)

func TestEncoding(t *testing.T) {
	data := "kind: pipeline\nname: pipeline\n---\nkind: secret\nname: secret\n---\nkind: signature\nhmac: signature"

	m, err := manifest.Decode(data)
	assert.NoError(t, err)
	assert.Equal(t, "pipeline", m.Resources[0].GetKind())
	assert.Equal(t, "secret", m.Resources[1].GetKind())
	assert.Equal(t, "signature", m.Resources[2].GetKind())

	encoded, err := manifest.Encode(m)
	assert.NoError(t, err)

	docs := strings.Split(encoded, "\n---")
	assert.Contains(t, docs[0], "kind: pipeline")
	assert.Contains(t, docs[1], "kind: secret")
	assert.Contains(t, docs[2], "kind: signature")
}
