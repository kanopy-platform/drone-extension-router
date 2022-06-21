package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/99designs/httpsignatures-go"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	"github.com/stretchr/testify/assert"
)

var testHandler http.Handler

func TestMain(m *testing.M) {
	testHandler = New("thisisnotsafe")
	os.Exit(m.Run())
}

func TestHandleRoot(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	pipeline := "name: default"
	req := newSignedRequest(t, converter.V1, converter.Request{Config: drone.Config{Data: pipeline}})

	testHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, fmt.Sprintf(`{"data":"%s","kind":""}`, pipeline), w.Body.String())
}

func TestHandleHealthz(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	testHandler.ServeHTTP(w, httptest.NewRequest("GET", "/healthz", nil))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	want := map[string]string{"status": "ok"}
	got := map[string]string{}

	assert.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, want, got)
}

// https://github.com/drone/drone-go/blob/v1.7.1/plugin/converter/handler_test.go#L31
func newSignedRequest(t *testing.T, accept string, extReq interface{}) *http.Request {
	// marshal extension request to json string
	reqBytes, err := json.Marshal(extReq)
	assert.NoError(t, err)

	// setup http request
	httpReq := httptest.NewRequest("GET", "/", bytes.NewReader(reqBytes))
	httpReq.Header.Add("Accept", accept)
	httpReq.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))

	// sign http request
	assert.NoError(t, httpsignatures.DefaultSha256Signer.AuthRequest("hmac-key", "thisisnotsafe", httpReq))

	return httpReq
}
