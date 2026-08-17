package fileserver

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type deadlineRoundTripper struct {
	requests int
	missing  bool
}

func (r *deadlineRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r.requests++
	if _, ok := req.Context().Deadline(); !ok {
		r.missing = true
	}
	body := "{}"
	if req.URL.Path == loginURI {
		body = `{"access_token":"test-token"}`
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}, nil
}

func TestFileServerRequestsHaveDeadlines(t *testing.T) {
	tests := []struct {
		name string
		call func(*SquirrelClient) error
	}{
		{name: "login", call: func(client *SquirrelClient) error { return client.Login() }},
		{name: "upload", call: func(client *SquirrelClient) error {
			return client.Upload("/images", "test.png", strings.NewReader("image"))
		}},
		{name: "health", call: func(client *SquirrelClient) error { return client.HealthCheck() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport := &deadlineRoundTripper{}
			client := &SquirrelClient{
				Host:       "http://fileserver.invalid",
				httpClient: &http.Client{Transport: transport},
			}
			if err := tt.call(client); err != nil {
				t.Fatalf("request error: %v", err)
			}
			if transport.requests != 1 || transport.missing {
				t.Fatalf("requests=%d missing_deadline=%t", transport.requests, transport.missing)
			}
		})
	}
}

func TestSharedHTTPClientHasTransportTimeouts(t *testing.T) {
	client := newHTTPClient()
	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("transport type = %T", client.Transport)
	}
	if client.Timeout <= 0 || transport.TLSHandshakeTimeout <= 0 || transport.ResponseHeaderTimeout <= 0 {
		t.Fatal("file server HTTP timeouts must all be bounded")
	}
}
