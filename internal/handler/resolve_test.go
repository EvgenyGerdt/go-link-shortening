package handler

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &mockURLService{
				resolveURL: tt.resolveURL,
			}

			handler := &Handler{
				service: service,
			}

			router := gin.New()
			router.GET("/:id", handler.ResolveURL)

			req := httptest.NewRequest(
				http.MethodGet,
				"/"+tt.id,
				nil,
			)

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					t.Fatal(err)
				}
			}(resp.Body)

			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf(
					"expected status %d, got %d",
					tt.expectedStatusCode,
					resp.StatusCode,
				)
			}

			if location := resp.Header.Get("Location"); location != tt.expectedURL {
				t.Errorf(
					"expected url %s, got %s",
					tt.expectedURL,
					location,
				)
			}
		})
	}
}
