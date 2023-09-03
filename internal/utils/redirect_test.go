package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirect(t *testing.T) {
	tests := []struct {
		name           string
		listenAddress  string
		httpsPort      string
		requestPath    string
		expectedStatus int
		expectedURL    string
	}{
		{
			name:           "standard port",
			listenAddress:  "example.com",
			httpsPort:      "443",
			requestPath:    "/somepath",
			expectedStatus: http.StatusMovedPermanently,
			expectedURL:    "https://example.com/somepath",
		},
		{
			name:           "non-standard port",
			listenAddress:  "example.com",
			httpsPort:      "8443",
			requestPath:    "/anotherpath",
			expectedStatus: http.StatusMovedPermanently,
			expectedURL:    "https://example.com:8443/anotherpath",
		},
		{
			name:           "empty path",
			listenAddress:  "example.com",
			httpsPort:      "443",
			requestPath:    "",
			expectedStatus: http.StatusMovedPermanently,
			expectedURL:    "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://"+tt.listenAddress+tt.requestPath, nil)
			w := httptest.NewRecorder()

			handler := Redirect(tt.listenAddress, tt.httpsPort)
			handler(w, req)

			resp := w.Result()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, resp.StatusCode)
			}

			location, ok := resp.Header["Location"]
			if !ok {
				t.Errorf("Location header not set")
			} else if location[0] != tt.expectedURL {
				t.Errorf("Expected URL %s; got %s", tt.expectedURL, location[0])
			}
		})
	}
}
