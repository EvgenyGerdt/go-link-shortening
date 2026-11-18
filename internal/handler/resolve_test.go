package handler

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ResolveURL(t *testing.T) {
	tests := []struct {
		name               string
		id                 string
		resolveURL         func(string) (string, error)
		expectedStatusCode int
		expectedURL        string
	}{
		{
			name: "missing id",
			id:   "",
			resolveURL: func(id string) (string, error) {
				t.Fatal("service should not be called")
				return "", nil
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "service error",
			id:   "test123",
			resolveURL: func(id string) (string, error) {
				return "", errors.New("URL not found")
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "success",
			id:   "test123",
			resolveURL: func(id string) (string, error) {
				return "http://example.com", nil
			},
			expectedStatusCode: http.StatusTemporaryRedirect,
			expectedURL:        "http://example.com",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := &mockURLService{
				resolveURL: test.resolveURL,
			}
			handler := &Handler{
				service: service,
			}
			req := httptest.NewRequest(http.MethodGet, "/"+test.id, nil)
			req.SetPathValue("id", test.id)

			recorder := httptest.NewRecorder()
			handler.ResolveURL(recorder, req)

			resp := recorder.Result()
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					t.Fatal(err)
				}
			}(resp.Body)
			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf("expected status %d, got %d", test.expectedStatusCode, resp.StatusCode)
			}

			if location := resp.Header.Get("Location"); location != test.expectedURL {
				t.Errorf("expected url %s, got %s", test.expectedURL, location)
			}
		})
	}
}
