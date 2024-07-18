package llm_test

import (
	"bytes"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ivanvanderbyl/graphrag-go/pkg/llm"
	"github.com/stretchr/testify/assert"
)

func TestCachingResponses(t *testing.T) {
	a := assert.New(t)

	tmpDir := t.TempDir()
	t.Log(tmpDir)

	// Create a listener on a fixed port
	listener, err := net.Listen("tcp", "localhost:54321")
	a.NoError(err)
	defer listener.Close()

	// Create a test server with the custom listener
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/success" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Success response"))
		} else if r.URL.Path == "/error" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error response"))
		}
	}))
	server.Listener = listener
	server.Start()
	defer server.Close()

	transport := llm.NewCacheTransport(http.DefaultTransport, nil, tmpDir, 0)

	tests := []struct {
		name           string
		req            *http.Request
		expected       int
		shouldBeCached bool
	}{
		{
			name:           "success",
			req:            httptest.NewRequest("POST", server.URL+"/success", bytes.NewBufferString("Completion Request")),
			expected:       http.StatusOK,
			shouldBeCached: true,
		},
		{
			name:           "success on read",
			req:            httptest.NewRequest("POST", server.URL+"/success", bytes.NewBufferString("Completion Request")),
			expected:       http.StatusOK,
			shouldBeCached: true,
		},
		{
			name:           "error",
			req:            httptest.NewRequest("POST", server.URL+"/error", nil),
			expected:       http.StatusInternalServerError,
			shouldBeCached: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			key, err := transport.GetCacheKey(tt.req)
			a.NoError(err)

			resp, err := transport.RoundTrip(tt.req)
			a.NoError(err)
			a.NotNil(resp)
			a.Equal(tt.expected, resp.StatusCode)

			_, err = os.Stat(filepath.Join(tmpDir, key))
			if tt.shouldBeCached {
				a.NoError(err)
			} else {
				a.Error(err)
			}
		})
	}
}
