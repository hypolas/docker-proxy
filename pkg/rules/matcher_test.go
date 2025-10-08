package rules

import (
	"testing"

	"dockershield/config"
)

func TestNewMatcher(t *testing.T) {
	rules := &config.AccessRules{
		Ping:       true,
		Containers: true,
		Post:       true,
	}

	matcher := NewMatcher(rules)
	if matcher == nil {
		t.Fatal("Expected non-nil matcher")
	}
	if matcher.rules != rules {
		t.Error("Expected matcher rules to match provided rules")
	}
}

func TestIsMethodAllowed(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		rules    *config.AccessRules
		expected bool
	}{
		{
			name:     "GET is always allowed",
			method:   "GET",
			rules:    &config.AccessRules{Post: false, Delete: false, Put: false},
			expected: true,
		},
		{
			name:     "HEAD is always allowed",
			method:   "HEAD",
			rules:    &config.AccessRules{Post: false, Delete: false, Put: false},
			expected: true,
		},
		{
			name:     "POST allowed when enabled",
			method:   "POST",
			rules:    &config.AccessRules{Post: true},
			expected: true,
		},
		{
			name:     "POST denied when disabled",
			method:   "POST",
			rules:    &config.AccessRules{Post: false},
			expected: false,
		},
		{
			name:     "DELETE allowed when enabled",
			method:   "DELETE",
			rules:    &config.AccessRules{Delete: true},
			expected: true,
		},
		{
			name:     "DELETE denied when disabled",
			method:   "DELETE",
			rules:    &config.AccessRules{Delete: false},
			expected: false,
		},
		{
			name:     "PUT allowed when enabled",
			method:   "PUT",
			rules:    &config.AccessRules{Put: true},
			expected: true,
		},
		{
			name:     "PUT denied when disabled",
			method:   "PUT",
			rules:    &config.AccessRules{Put: false},
			expected: false,
		},
		{
			name:     "PATCH allowed when enabled",
			method:   "PATCH",
			rules:    &config.AccessRules{Put: true},
			expected: true,
		},
		{
			name:     "Unknown method denied",
			method:   "TRACE",
			rules:    &config.AccessRules{Post: true, Delete: true, Put: true},
			expected: false,
		},
		{
			name:     "Lowercase method GET",
			method:   "get",
			rules:    &config.AccessRules{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := NewMatcher(tt.rules)
			result := matcher.isMethodAllowed(tt.method)
			if result != tt.expected {
				t.Errorf("Expected %v for method %s, got %v", tt.expected, tt.method, result)
			}
		})
	}
}

func TestIsPathAllowed(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		rules    *config.AccessRules
		expected bool
	}{
		{
			name:     "Ping endpoint allowed",
			path:     "/v1.41/_ping",
			rules:    &config.AccessRules{Ping: true},
			expected: true,
		},
		{
			name:     "Ping endpoint denied",
			path:     "/v1.41/_ping",
			rules:    &config.AccessRules{Ping: false},
			expected: false,
		},
		{
			name:     "Containers endpoint allowed",
			path:     "/v1.41/containers/json",
			rules:    &config.AccessRules{Containers: true},
			expected: true,
		},
		{
			name:     "Containers endpoint denied",
			path:     "/v1.41/containers/json",
			rules:    &config.AccessRules{Containers: false},
			expected: false,
		},
		{
			name:     "Images endpoint allowed",
			path:     "/v1.41/images/json",
			rules:    &config.AccessRules{Images: true},
			expected: true,
		},
		{
			name:     "Images endpoint denied",
			path:     "/v1.41/images/json",
			rules:    &config.AccessRules{Images: false},
			expected: false,
		},
		{
			name:     "Version endpoint allowed",
			path:     "/version",
			rules:    &config.AccessRules{Version: true},
			expected: true,
		},
		{
			name:     "Events endpoint allowed",
			path:     "/v1.43/events",
			rules:    &config.AccessRules{Events: true},
			expected: true,
		},
		{
			name:     "Build endpoint allowed",
			path:     "/v1.41/build",
			rules:    &config.AccessRules{Build: true},
			expected: true,
		},
		{
			name:     "Networks endpoint allowed",
			path:     "/v1.41/networks/create",
			rules:    &config.AccessRules{Networks: true},
			expected: true,
		},
		{
			name:     "Volumes endpoint allowed",
			path:     "/v1.41/volumes",
			rules:    &config.AccessRules{Volumes: true},
			expected: true,
		},
		{
			name:     "Exec endpoint allowed",
			path:     "/v1.41/exec/abc123/start",
			rules:    &config.AccessRules{Exec: true},
			expected: true,
		},
		{
			name:     "Unknown endpoint denied",
			path:     "/v1.41/unknown/endpoint",
			rules:    &config.AccessRules{},
			expected: false,
		},
		{
			name:     "Path without version",
			path:     "/containers/json",
			rules:    &config.AccessRules{Containers: true},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := NewMatcher(tt.rules)
			result := matcher.isPathAllowed(tt.path)
			if result != tt.expected {
				t.Errorf("Expected %v for path %s, got %v", tt.expected, tt.path, result)
			}
		})
	}
}

func TestIsAllowed(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		path     string
		rules    *config.AccessRules
		expected bool
	}{
		{
			name:   "GET containers allowed",
			method: "GET",
			path:   "/v1.41/containers/json",
			rules: &config.AccessRules{
				Containers: true,
			},
			expected: true,
		},
		{
			name:   "POST containers allowed when POST and Containers enabled",
			method: "POST",
			path:   "/v1.41/containers/create",
			rules: &config.AccessRules{
				Containers: true,
				Post:       true,
			},
			expected: true,
		},
		{
			name:   "POST containers denied when POST disabled",
			method: "POST",
			path:   "/v1.41/containers/create",
			rules: &config.AccessRules{
				Containers: true,
				Post:       false,
			},
			expected: false,
		},
		{
			name:   "POST containers denied when Containers disabled",
			method: "POST",
			path:   "/v1.41/containers/create",
			rules: &config.AccessRules{
				Containers: false,
				Post:       true,
			},
			expected: false,
		},
		{
			name:   "DELETE images denied when DELETE disabled",
			method: "DELETE",
			path:   "/v1.41/images/nginx",
			rules: &config.AccessRules{
				Images: true,
				Delete: false,
			},
			expected: false,
		},
		{
			name:   "DELETE images allowed when both enabled",
			method: "DELETE",
			path:   "/v1.41/images/nginx",
			rules: &config.AccessRules{
				Images: true,
				Delete: true,
			},
			expected: true,
		},
		{
			name:   "GET always allowed for enabled endpoint",
			method: "GET",
			path:   "/v1.41/images/json",
			rules: &config.AccessRules{
				Images: true,
				Post:   false,
				Delete: false,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := NewMatcher(tt.rules)
			result := matcher.IsAllowed(tt.method, tt.path)
			if result != tt.expected {
				t.Errorf("Expected %v for %s %s, got %v", tt.expected, tt.method, tt.path, result)
			}
		})
	}
}

func TestRemoveAPIVersion(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Remove v1.41",
			path:     "/v1.41/containers/json",
			expected: "/containers/json",
		},
		{
			name:     "Remove v1.43",
			path:     "/v1.43/images/json",
			expected: "/images/json",
		},
		{
			name:     "No version to remove",
			path:     "/containers/json",
			expected: "/containers/json",
		},
		{
			name:     "Version at end",
			path:     "/containers/v1.41",
			expected: "/containers/v1.41",
		},
		{
			name:     "Multiple version-like patterns",
			path:     "/v1.41/containers/v1.42/json",
			expected: "/containers/v1.42/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeAPIVersion(tt.path)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
