package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"dockershield/config"
	"dockershield/pkg/rules"

	"github.com/gin-gonic/gin"
)

func TestACLMiddleware(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		path           string
		rules          *config.AccessRules
		expectedStatus int
		expectAborted  bool
	}{
		{
			name:   "GET containers allowed",
			method: "GET",
			path:   "/v1.41/containers/json",
			rules: &config.AccessRules{
				Containers: true,
			},
			expectedStatus: http.StatusOK,
			expectAborted:  false,
		},
		{
			name:   "GET containers denied",
			method: "GET",
			path:   "/v1.41/containers/json",
			rules: &config.AccessRules{
				Containers: false,
			},
			expectedStatus: http.StatusForbidden,
			expectAborted:  true,
		},
		{
			name:   "POST containers allowed",
			method: "POST",
			path:   "/v1.41/containers/create",
			rules: &config.AccessRules{
				Containers: true,
				Post:       true,
			},
			expectedStatus: http.StatusOK,
			expectAborted:  false,
		},
		{
			name:   "POST containers denied - endpoint disabled",
			method: "POST",
			path:   "/v1.41/containers/create",
			rules: &config.AccessRules{
				Containers: false,
				Post:       true,
			},
			expectedStatus: http.StatusForbidden,
			expectAborted:  true,
		},
		{
			name:   "POST containers denied - method disabled",
			method: "POST",
			path:   "/v1.41/containers/create",
			rules: &config.AccessRules{
				Containers: true,
				Post:       false,
			},
			expectedStatus: http.StatusForbidden,
			expectAborted:  true,
		},
		{
			name:   "DELETE images allowed",
			method: "DELETE",
			path:   "/v1.41/images/nginx",
			rules: &config.AccessRules{
				Images: true,
				Delete: true,
			},
			expectedStatus: http.StatusOK,
			expectAborted:  false,
		},
		{
			name:   "DELETE images denied",
			method: "DELETE",
			path:   "/v1.41/images/nginx",
			rules: &config.AccessRules{
				Images: true,
				Delete: false,
			},
			expectedStatus: http.StatusForbidden,
			expectAborted:  true,
		},
		{
			name:   "Ping always allowed with GET",
			method: "GET",
			path:   "/v1.41/_ping",
			rules: &config.AccessRules{
				Ping: true,
			},
			expectedStatus: http.StatusOK,
			expectAborted:  false,
		},
		{
			name:   "Version always allowed with GET",
			method: "GET",
			path:   "/version",
			rules: &config.AccessRules{
				Version: true,
			},
			expectedStatus: http.StatusOK,
			expectAborted:  false,
		},
		{
			name:   "Unknown endpoint denied",
			method: "GET",
			path:   "/v1.41/unknown",
			rules: &config.AccessRules{
				Containers: true,
				Images:     true,
			},
			expectedStatus: http.StatusForbidden,
			expectAborted:  true,
		},
		{
			name:   "HEAD method allowed",
			method: "HEAD",
			path:   "/v1.41/containers/json",
			rules: &config.AccessRules{
				Containers: true,
			},
			expectedStatus: http.StatusOK,
			expectAborted:  false,
		},
		{
			name:   "PUT method denied",
			method: "PUT",
			path:   "/v1.41/containers/update",
			rules: &config.AccessRules{
				Containers: true,
				Put:        false,
			},
			expectedStatus: http.StatusForbidden,
			expectAborted:  true,
		},
		{
			name:   "PUT method allowed",
			method: "PUT",
			path:   "/v1.41/containers/update",
			rules: &config.AccessRules{
				Containers: true,
				Put:        true,
			},
			expectedStatus: http.StatusOK,
			expectAborted:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create matcher
			matcher := rules.NewMatcher(tt.rules)

			// Create router
			router := gin.New()
			router.Use(ACLMiddleware(matcher))

			// Add a catch-all route for testing
			router.Any("/*path", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			// Create request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			// Serve request
			router.ServeHTTP(w, req)

			// Check status
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response for forbidden requests
			if tt.expectAborted && w.Code == http.StatusForbidden {
				// Response should contain JSON error message
				if w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
					t.Error("Expected JSON content type for forbidden response")
				}
			}
		})
	}
}

func TestACLMiddlewareJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accessRules := &config.AccessRules{
		Containers: false,
	}
	matcher := rules.NewMatcher(accessRules)

	router := gin.New()
	router.Use(ACLMiddleware(matcher))
	router.Any("/*path", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/v1.41/containers/json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}

	// Check that response contains expected fields
	body := w.Body.String()
	if body == "" {
		t.Error("Expected non-empty response body")
	}

	// Should contain JSON with message, path, and method
	expectedFields := []string{`"message"`, `"path"`, `"method"`}
	for _, field := range expectedFields {
		if !contains(body, field) {
			t.Errorf("Expected response to contain %s", field)
		}
	}
}

func TestACLMiddlewareAbort(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accessRules := &config.AccessRules{
		Containers: false,
	}
	matcher := rules.NewMatcher(accessRules)

	handlerCalled := false

	router := gin.New()
	router.Use(ACLMiddleware(matcher))
	router.GET("/v1.41/containers/json", func(c *gin.Context) {
		handlerCalled = true
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/v1.41/containers/json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if handlerCalled {
		t.Error("Handler should not have been called when request is forbidden")
	}

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestACLMiddlewareNext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accessRules := &config.AccessRules{
		Containers: true,
	}
	matcher := rules.NewMatcher(accessRules)

	handlerCalled := false

	router := gin.New()
	router.Use(ACLMiddleware(matcher))
	router.GET("/v1.41/containers/json", func(c *gin.Context) {
		handlerCalled = true
		c.JSON(http.StatusOK, gin.H{"result": "ok"})
	})

	req := httptest.NewRequest("GET", "/v1.41/containers/json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if !handlerCalled {
		t.Error("Handler should have been called when request is allowed")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestACLMiddlewareMultipleEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accessRules := &config.AccessRules{
		Containers: true,
		Images:     false,
		Networks:   true,
		Volumes:    false,
	}
	matcher := rules.NewMatcher(accessRules)

	router := gin.New()
	router.Use(ACLMiddleware(matcher))
	router.Any("/*path", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	tests := []struct {
		path           string
		expectedStatus int
	}{
		{"/v1.41/containers/json", http.StatusOK},
		{"/v1.41/images/json", http.StatusForbidden},
		{"/v1.41/networks", http.StatusOK},
		{"/v1.41/volumes", http.StatusForbidden},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", tt.path, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != tt.expectedStatus {
			t.Errorf("For path %s: expected status %d, got %d", tt.path, tt.expectedStatus, w.Code)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && hasSubstring(s, substr))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
