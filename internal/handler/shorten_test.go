package handler

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandler_ShortenURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		body               string
		shortenURL         func(string) (string, error)
		expectedStatusCode int
		expectedBody       string
		serviceCalled      bool
	}{
		{
			name:               "empty body",
			body:               "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "invalid URL",
			body:               "://bad",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "missing scheme",
			body:               "example.com",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "unsupported scheme",
			body:               "ftp://example.com",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "missing host",
			body:               "https://",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: "https://example.com",
			shortenURL: func(url string) (string, error) {
				return "", errors.New("service error")
			},
			expectedStatusCode: http.StatusBadRequest,
			serviceCalled:      true,
		},
		{
			name: "success",
			body: "https://example.com",
			shortenURL: func(url string) (string, error) {
				return "abc123", nil
			},
			expectedStatusCode: http.StatusCreated,
			expectedBody:       "http://short.test/abc123",
			serviceCalled:      true,
		},
		{
			name: "success with query params",
			body: "https://example.com/path?q=1",
			shortenURL: func(url string) (string, error) {
				return "xyz789", nil
			},
			expectedStatusCode: http.StatusCreated,
			expectedBody:       "http://short.test/xyz789",
			serviceCalled:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false

			service := &mockURLService{
				shortenURL: func(url string) (string, error) {
					called = true

					if url != tt.body {
						t.Errorf(
							"expected service URL %q, got %q",
							tt.body,
							url,
						)
					}

					if tt.shortenURL == nil {
						t.Fatal("unexpected service call")
						return "", nil
					}

					return tt.shortenURL(url)
				},
			}

			handler := &Handler{
				service: service,
			}

			router := gin.New()
			router.POST("/", handler.ShortenURL)

			// Создаём HTTP-запрос
			req := httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(tt.body),
			)

			req.Host = "short.test"
			req.Header.Set("Content-Type", "text/plain")

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

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if string(body) != tt.expectedBody {
				t.Errorf(
					"expected body %q, got %q",
					tt.expectedBody,
					string(body),
				)
			}

			if tt.expectedStatusCode == http.StatusCreated {
				contentType := resp.Header.Get("Content-Type")

				mediaType, _, err := mime.ParseMediaType(contentType)
				if err != nil {
					t.Fatalf("invalid Content-Type %q: %v", contentType, err)
				}

				if mediaType != "text/plain" {
					t.Errorf(
						"expected Content-Type text/plain, got %q",
						mediaType,
					)
				}
			}
		})
	}
}
