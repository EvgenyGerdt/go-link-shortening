package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandler_ResolveURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		id                 string
		resolveURL         func(string) (string, error)
		expectedStatusCode int
		expectedURL        string
		serviceCalled      bool
	}{
		{
			name:               "missing id",
			id:                 "",
			expectedStatusCode: http.StatusNotFound,
			serviceCalled:      false,
		},
		{
			name: "service error",
			id:   "test123",
			resolveURL: func(id string) (string, error) {
				return "", errors.New("URL not found")
			},
			expectedStatusCode: http.StatusBadRequest,
			serviceCalled:      true,
		},
		{
			name: "success",
			id:   "test123",
			resolveURL: func(id string) (string, error) {
				return "http://example.com", nil
			},
			expectedStatusCode: http.StatusTemporaryRedirect,
			expectedURL:        "http://example.com",
			serviceCalled:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false

			service := &mockURLService{
				resolveURL: func(id string) (string, error) {
					called = true

					if id != tt.id {
						t.Errorf(
							"expected id %q, got %q",
							tt.id,
							id,
						)
					}

					if tt.resolveURL == nil {
						t.Fatal("unexpected service call")
						return "", nil
					}

					return tt.resolveURL(id)
				},
			}

			handler := &Handler{
				service: service,
				baseURL: "http://short.test",
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
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf(
					"expected status %d, got %d",
					tt.expectedStatusCode,
					resp.StatusCode,
				)
			}

			if called != tt.serviceCalled {
				t.Errorf(
					"expected service called %t, got %t",
					tt.serviceCalled,
					called,
				)
			}

			if location := resp.Header.Get("Location"); location != tt.expectedURL {
				t.Errorf(
					"expected Location %q, got %q",
					tt.expectedURL,
					location,
				)
			}
		})
	}
}
