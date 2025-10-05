package rules

import (
	"regexp"
	"strings"

	"docker-proxy/config"
)

// Matcher determines if a request is allowed based on access rules
type Matcher struct {
	rules *config.AccessRules
}

// NewMatcher creates a new rule matcher
func NewMatcher(rules *config.AccessRules) *Matcher {
	return &Matcher{rules: rules}
}

// IsAllowed checks if a request with given method and path is allowed
func (m *Matcher) IsAllowed(method, path string) bool {
	// Check HTTP method first
	if !m.isMethodAllowed(method) {
		return false
	}

	// Check API endpoint
	return m.isPathAllowed(path)
}

// isMethodAllowed checks if the HTTP method is allowed
func (m *Matcher) isMethodAllowed(method string) bool {
	method = strings.ToUpper(method)

	switch method {
	case "GET", "HEAD":
		return true // Always allow read operations
	case "POST":
		return m.rules.Post
	case "DELETE":
		return m.rules.Delete
	case "PUT", "PATCH":
		return m.rules.Put
	default:
		return false
	}
}

// isPathAllowed checks if the API path is allowed
func (m *Matcher) isPathAllowed(path string) bool {
	// Remove API version prefix
	path = removeAPIVersion(path)

	// Check against each endpoint pattern
	patterns := map[*regexp.Regexp]bool{
		regexp.MustCompile(`^/_ping`):                  m.rules.Ping,
		regexp.MustCompile(`^/events`):                 m.rules.Events,
		regexp.MustCompile(`^/version`):                m.rules.Version,
		regexp.MustCompile(`^/auth`):                   m.rules.Auth,
		regexp.MustCompile(`^/build`):                  m.rules.Build,
		regexp.MustCompile(`^/commit`):                 m.rules.Commit,
		regexp.MustCompile(`^/configs`):                m.rules.Configs,
		regexp.MustCompile(`^/containers`):             m.rules.Containers,
		regexp.MustCompile(`^/distribution`):           m.rules.Distribution,
		regexp.MustCompile(`^/exec`):                   m.rules.Exec,
		regexp.MustCompile(`^/images`):                 m.rules.Images,
		regexp.MustCompile(`^/info`):                   m.rules.Info,
		regexp.MustCompile(`^/networks`):               m.rules.Networks,
		regexp.MustCompile(`^/nodes`):                  m.rules.Nodes,
		regexp.MustCompile(`^/plugins`):                m.rules.Plugins,
		regexp.MustCompile(`^/secrets`):                m.rules.Secrets,
		regexp.MustCompile(`^/services`):               m.rules.Services,
		regexp.MustCompile(`^/session`):                m.rules.Session,
		regexp.MustCompile(`^/swarm`):                  m.rules.Swarm,
		regexp.MustCompile(`^/system`):                 m.rules.System,
		regexp.MustCompile(`^/tasks`):                  m.rules.Tasks,
		regexp.MustCompile(`^/volumes`):                m.rules.Volumes,
	}

	for pattern, allowed := range patterns {
		if pattern.MatchString(path) {
			return allowed
		}
	}

	// Default deny unknown endpoints
	return false
}

// removeAPIVersion removes the API version prefix from the path
func removeAPIVersion(path string) string {
	// Match patterns like /v1.41/containers
	re := regexp.MustCompile(`^/v\d+\.\d+`)
	return re.ReplaceAllString(path, "")
}
