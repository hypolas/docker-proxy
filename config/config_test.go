package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "Environment variable set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
		{
			name:         "Environment variable not set",
			key:          "TEST_VAR_UNSET",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "Empty default value",
			key:          "TEST_VAR_EMPTY",
			defaultValue: "",
			envValue:     "",
			expected:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestGetBoolEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue bool
		envValue     string
		expected     bool
	}{
		{
			name:         "Value '1' returns true",
			key:          "TEST_BOOL",
			defaultValue: false,
			envValue:     "1",
			expected:     true,
		},
		{
			name:         "Value '0' returns false",
			key:          "TEST_BOOL",
			defaultValue: true,
			envValue:     "0",
			expected:     false,
		},
		{
			name:         "Value 'true' returns true",
			key:          "TEST_BOOL",
			defaultValue: false,
			envValue:     "true",
			expected:     true,
		},
		{
			name:         "Value 'True' returns true",
			key:          "TEST_BOOL",
			defaultValue: false,
			envValue:     "True",
			expected:     true,
		},
		{
			name:         "Value 'false' returns false",
			key:          "TEST_BOOL",
			defaultValue: true,
			envValue:     "false",
			expected:     false,
		},
		{
			name:         "Value 'False' returns false",
			key:          "TEST_BOOL",
			defaultValue: true,
			envValue:     "False",
			expected:     false,
		},
		{
			name:         "Value 'yes' returns true",
			key:          "TEST_BOOL",
			defaultValue: false,
			envValue:     "yes",
			expected:     true,
		},
		{
			name:         "Value 'no' returns false",
			key:          "TEST_BOOL",
			defaultValue: true,
			envValue:     "no",
			expected:     false,
		},
		{
			name:         "Empty value returns default",
			key:          "TEST_BOOL_EMPTY",
			defaultValue: true,
			envValue:     "",
			expected:     true,
		},
		{
			name:         "Invalid value returns default",
			key:          "TEST_BOOL",
			defaultValue: true,
			envValue:     "invalid",
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			} else {
				os.Unsetenv(tt.key)
			}
			defer os.Unsetenv(tt.key)

			result := getBoolEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestLoadAccessRules(t *testing.T) {
	// Clean environment
	cleanEnv := func() {
		envVars := []string{
			"EVENTS", "PING", "VERSION", "AUTH", "BUILD", "COMMIT",
			"CONFIGS", "CONTAINERS", "DISTRIBUTION", "EXEC", "IMAGES",
			"INFO", "NETWORKS", "NODES", "PLUGINS", "SECRETS", "SERVICES",
			"SESSION", "SWARM", "SYSTEM", "TASKS", "VOLUMES",
			"POST", "DELETE", "PUT",
		}
		for _, v := range envVars {
			os.Unsetenv(v)
		}
	}

	t.Run("Default values", func(t *testing.T) {
		cleanEnv()
		defer cleanEnv()

		rules := loadAccessRules()

		// Check defaults granted
		if !rules.Events {
			t.Error("Expected Events to be true by default")
		}
		if !rules.Ping {
			t.Error("Expected Ping to be true by default")
		}
		if !rules.Version {
			t.Error("Expected Version to be true by default")
		}

		// Check defaults denied
		if rules.Containers {
			t.Error("Expected Containers to be false by default")
		}
		if rules.Images {
			t.Error("Expected Images to be false by default")
		}
		if rules.Post {
			t.Error("Expected Post to be false by default")
		}
		if rules.Delete {
			t.Error("Expected Delete to be false by default")
		}
		if rules.Put {
			t.Error("Expected Put to be false by default")
		}
	})

	t.Run("Custom values", func(t *testing.T) {
		cleanEnv()
		defer cleanEnv()

		os.Setenv("CONTAINERS", "1")
		os.Setenv("IMAGES", "1")
		os.Setenv("POST", "1")
		os.Setenv("PING", "0")

		rules := loadAccessRules()

		if !rules.Containers {
			t.Error("Expected Containers to be true")
		}
		if !rules.Images {
			t.Error("Expected Images to be true")
		}
		if !rules.Post {
			t.Error("Expected Post to be true")
		}
		if rules.Ping {
			t.Error("Expected Ping to be false")
		}
	})
}

func TestLoad(t *testing.T) {
	// Clean environment
	cleanEnv := func() {
		vars := []string{
			"LISTEN_ADDR", "LISTEN_SOCKET", "DOCKER_SOCKET",
			"LOG_LEVEL", "API_VERSION", "FILTERS_CONFIG",
		}
		for _, v := range vars {
			os.Unsetenv(v)
		}
	}

	t.Run("Default configuration", func(t *testing.T) {
		cleanEnv()
		defer cleanEnv()

		cfg := Load()

		if cfg.ListenAddr != ":2375" {
			t.Errorf("Expected ListenAddr ':2375', got '%s'", cfg.ListenAddr)
		}
		if cfg.ListenSocket != "" {
			t.Errorf("Expected empty ListenSocket, got '%s'", cfg.ListenSocket)
		}
		if cfg.DockerSocket != "unix:///var/run/docker.sock" {
			t.Errorf("Expected DockerSocket 'unix:///var/run/docker.sock', got '%s'", cfg.DockerSocket)
		}
		if cfg.LogLevel != "info" {
			t.Errorf("Expected LogLevel 'info', got '%s'", cfg.LogLevel)
		}
		if cfg.APIVersion != "" {
			t.Errorf("Expected empty APIVersion, got '%s'", cfg.APIVersion)
		}
		if cfg.AccessRules == nil {
			t.Error("Expected non-nil AccessRules")
		}
	})

	t.Run("Custom configuration", func(t *testing.T) {
		cleanEnv()
		defer cleanEnv()

		os.Setenv("LISTEN_ADDR", ":3000")
		os.Setenv("LISTEN_SOCKET", "unix:///tmp/dockershield.sock")
		os.Setenv("DOCKER_SOCKET", "tcp://localhost:2376")
		os.Setenv("LOG_LEVEL", "debug")
		os.Setenv("API_VERSION", "1.41")

		cfg := Load()

		if cfg.ListenAddr != ":3000" {
			t.Errorf("Expected ListenAddr ':3000', got '%s'", cfg.ListenAddr)
		}
		if cfg.ListenSocket != "unix:///tmp/dockershield.sock" {
			t.Errorf("Expected ListenSocket 'unix:///tmp/dockershield.sock', got '%s'", cfg.ListenSocket)
		}
		if cfg.DockerSocket != "tcp://localhost:2376" {
			t.Errorf("Expected DockerSocket 'tcp://localhost:2376', got '%s'", cfg.DockerSocket)
		}
		if cfg.LogLevel != "debug" {
			t.Errorf("Expected LogLevel 'debug', got '%s'", cfg.LogLevel)
		}
		if cfg.APIVersion != "1.41" {
			t.Errorf("Expected APIVersion '1.41', got '%s'", cfg.APIVersion)
		}
	})

	t.Run("Unix socket formats", func(t *testing.T) {
		cleanEnv()
		defer cleanEnv()

		tests := []struct {
			name     string
			envValue string
			expected string
		}{
			{
				name:     "unix:// prefix",
				envValue: "unix:///tmp/docker.sock",
				expected: "unix:///tmp/docker.sock",
			},
			{
				name:     "Absolute path",
				envValue: "/tmp/docker.sock",
				expected: "/tmp/docker.sock",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				os.Setenv("LISTEN_SOCKET", tt.envValue)
				cfg := Load()
				if cfg.ListenSocket != tt.expected {
					t.Errorf("Expected '%s', got '%s'", tt.expected, cfg.ListenSocket)
				}
				os.Unsetenv("LISTEN_SOCKET")
			})
		}
	})
}

func TestAccessRulesDefaults(t *testing.T) {
	cleanEnv := func() {
		envVars := []string{
			"EVENTS", "PING", "VERSION", "AUTH", "BUILD", "COMMIT",
			"CONFIGS", "CONTAINERS", "DISTRIBUTION", "EXEC", "IMAGES",
			"INFO", "NETWORKS", "NODES", "PLUGINS", "SECRETS", "SERVICES",
			"SESSION", "SWARM", "SYSTEM", "TASKS", "VOLUMES",
			"POST", "DELETE", "PUT",
		}
		for _, v := range envVars {
			os.Unsetenv(v)
		}
	}

	cleanEnv()
	defer cleanEnv()

	rules := loadAccessRules()

	// Test default granted endpoints
	defaultGranted := map[string]bool{
		"Events":  rules.Events,
		"Ping":    rules.Ping,
		"Version": rules.Version,
	}

	for name, value := range defaultGranted {
		if !value {
			t.Errorf("Expected %s to be granted by default", name)
		}
	}

	// Test default denied endpoints
	defaultDenied := []struct {
		name  string
		value bool
	}{
		{"Auth", rules.Auth},
		{"Build", rules.Build},
		{"Commit", rules.Commit},
		{"Configs", rules.Configs},
		{"Containers", rules.Containers},
		{"Distribution", rules.Distribution},
		{"Exec", rules.Exec},
		{"Images", rules.Images},
		{"Info", rules.Info},
		{"Networks", rules.Networks},
		{"Nodes", rules.Nodes},
		{"Plugins", rules.Plugins},
		{"Secrets", rules.Secrets},
		{"Services", rules.Services},
		{"Session", rules.Session},
		{"Swarm", rules.Swarm},
		{"System", rules.System},
		{"Tasks", rules.Tasks},
		{"Volumes", rules.Volumes},
		{"Post", rules.Post},
		{"Delete", rules.Delete},
		{"Put", rules.Put},
	}

	for _, item := range defaultDenied {
		if item.value {
			t.Errorf("Expected %s to be denied by default", item.name)
		}
	}
}
