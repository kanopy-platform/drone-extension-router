package manifest_test

import (
	"testing"

	"github.com/kanopy-platform/drone-extension-router/pkg/manifest"
	"github.com/stretchr/testify/assert"
)

func TestDecodeEncode(t *testing.T) {
	tests := []struct {
		desc    string
		input   string
		decoded []manifest.Resource
		encoded string
	}{
		{
			desc: "test empty input",
		},
		{
			desc:    "test separator prefix",
			input:   "---\nkind: pipeline\n",
			decoded: []manifest.Resource{&manifest.Pipeline{Kind: manifest.KindPipeline}},
			encoded: "kind: pipeline",
		},
		{
			desc: "test multiple resources",
			input: `
kind: pipeline
name: pipeline
inline: pipeline
---
kind: secret
name: secret
inline: secret
---
kind: signature
hmac: signature
inline: signature
`,
			decoded: []manifest.Resource{
				&manifest.Pipeline{
					Kind: manifest.KindPipeline,
					Name: "pipeline",
					ResourceData: map[string]interface{}{
						"inline": "pipeline",
					},
				},
				&manifest.Secret{
					Kind: manifest.KindSecret,
					Name: "secret",
					ResourceData: map[string]interface{}{
						"inline": "secret",
					},
				},
				&manifest.Object{
					Kind: manifest.Kind("signature"),
					ResourceData: map[string]interface{}{
						"hmac":   "signature",
						"inline": "signature",
					},
				},
			},
			encoded: `kind: pipeline
name: pipeline
inline: pipeline
---
kind: secret
name: secret
inline: secret
---
kind: signature
hmac: signature
inline: signature`,
		},
	}

	for _, test := range tests {
		gotDecoded, err := manifest.Decode(test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.decoded, gotDecoded, "Decoding: "+test.desc)

		gotEncoded, err := manifest.Encode(test.decoded)
		assert.NoError(t, err)
		assert.Equal(t, test.encoded, gotEncoded, "Encoding: "+test.desc)
	}
}

func BenchmarkDecode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = manifest.Decode("---\ntest: test")
	}
}

func BenchmarkEncode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = manifest.Encode([]manifest.Resource{&manifest.Pipeline{Kind: manifest.KindPipeline}})
	}
}
