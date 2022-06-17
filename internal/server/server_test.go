package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testHandler http.Handler

func TestMain(m *testing.M) {
	testHandler = New()
	os.Exit(m.Run())
}

func TestHandleRoot(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	testHandler.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	assert.Equal(t, http.StatusOK, w.Code)
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
