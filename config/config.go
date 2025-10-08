package config

import (
	"os"
	"strconv"
	"strings"

	"dockershield/pkg/filters"
)

// Config holds the application configuration
type Config struct {
	ListenAddr      string
	ListenSocket    string // Unix socket path for listening (optional)
	DockerSocket    string
	LogLevel        string
	APIVersion      string
	AccessRules     *AccessRules
	AdvancedFilters *filters.AdvancedFilter // Filtres avancés (optionnel)
	FiltersPath     string                  // Chemin vers le fichier JSON de filtres
}

// AccessRules defines which Docker API endpoints are allowed
type AccessRules struct {
	// Default granted
	Events  bool
	Ping    bool
	Version bool

	// API endpoints
	Auth         bool
	Build        bool
	Commit       bool
	Configs      bool
	Containers   bool
	Distribution bool
	Exec         bool
	Images       bool
	Info         bool
	Networks     bool
	Nodes        bool
	Plugins      bool
	Secrets      bool
	Services     bool
	Session      bool
	Swarm        bool
	System       bool
	Tasks        bool
	Volumes      bool

	// HTTP Methods
	Post   bool
	Delete bool
	Put    bool
}

// Load loads configuration from environment variables
func Load() *Config {
	filtersPath := getEnv("FILTERS_CONFIG", "")

	// Charger les filtres depuis JSON (si configuré)
	jsonFilters := loadAdvancedFilters(filtersPath)

	// Charger les filtres depuis les variables d'environnement (prioritaire)
	envFilters := LoadFiltersFromEnv()

	// Fusionner avec priorité aux env vars
	mergedFilters := MergeFilters(jsonFilters, envFilters)

	// Appliquer les filtres par défaut si l'utilisateur n'a pas désactivé les défauts
	if !CanOverrideDefaults() {
		mergedFilters = ApplyDefaults(mergedFilters)
	}

	config := &Config{
		ListenAddr:      getEnv("LISTEN_ADDR", ":2375"),
		ListenSocket:    getEnv("LISTEN_SOCKET", ""),
		DockerSocket:    getEnv("DOCKER_SOCKET", "unix:///var/run/docker.sock"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		APIVersion:      getEnv("API_VERSION", ""), // Will be auto-detected if empty
		AccessRules:     loadAccessRules(),
		FiltersPath:     filtersPath,
		AdvancedFilters: mergedFilters,
	}
	return config
}

// loadAccessRules loads access rules from environment variables
func loadAccessRules() *AccessRules {
	return &AccessRules{
		// Default granted
		Events:  getBoolEnv("EVENTS", true),
		Ping:    getBoolEnv("PING", true),
		Version: getBoolEnv("VERSION", true),

		// API endpoints - default denied
		Auth:         getBoolEnv("AUTH", false),
		Build:        getBoolEnv("BUILD", false),
		Commit:       getBoolEnv("COMMIT", false),
		Configs:      getBoolEnv("CONFIGS", false),
		Containers:   getBoolEnv("CONTAINERS", false),
		Distribution: getBoolEnv("DISTRIBUTION", false),
		Exec:         getBoolEnv("EXEC", false),
		Images:       getBoolEnv("IMAGES", false),
		Info:         getBoolEnv("INFO", false),
		Networks:     getBoolEnv("NETWORKS", false),
		Nodes:        getBoolEnv("NODES", false),
		Plugins:      getBoolEnv("PLUGINS", false),
		Secrets:      getBoolEnv("SECRETS", false),
		Services:     getBoolEnv("SERVICES", false),
		Session:      getBoolEnv("SESSION", false),
		Swarm:        getBoolEnv("SWARM", false),
		System:       getBoolEnv("SYSTEM", false),
		Tasks:        getBoolEnv("TASKS", false),
		Volumes:      getBoolEnv("VOLUMES", false),

		// HTTP Methods - POST, DELETE, PUT denied by default (read-only mode)
		Post:   getBoolEnv("POST", false),
		Delete: getBoolEnv("DELETE", false),
		Put:    getBoolEnv("PUT", false),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getBoolEnv gets a boolean environment variable or returns a default value
func getBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	// Support "0" and "1" as per docker-socket-proxy
	value = strings.ToLower(value)
	if value == "1" || value == "true" || value == "yes" {
		return true
	}
	if value == "0" || value == "false" || value == "no" {
		return false
	}

	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolVal
}

// loadAdvancedFilters loads advanced filters from JSON file
func loadAdvancedFilters(filtersPath string) *filters.AdvancedFilter {
	if filtersPath == "" {
		return nil
	}

	data, err := os.ReadFile(filtersPath)
	if err != nil {
		return nil
	}

	filter, err := filters.LoadFromJSON(data)
	if err != nil {
		return nil
	}

	return filter
}
