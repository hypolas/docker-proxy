package config

import (
	"os"
	"strings"

	"docker-proxy/pkg/filters"
)

const envPrefix = "DKRPRX__"

// LoadFiltersFromEnv charge les filtres depuis les variables d'environnement
// Les variables d'environnement sont prioritaires sur le fichier JSON
func LoadFiltersFromEnv() *filters.AdvancedFilter {
	filter := &filters.AdvancedFilter{}
	hasAnyFilter := false

	// Volumes
	if vf := loadVolumeFilters(); vf != nil {
		filter.Volumes = vf
		hasAnyFilter = true
	}

	// Containers
	if cf := loadContainerFilters(); cf != nil {
		filter.Containers = cf
		hasAnyFilter = true
	}

	// Networks
	if nf := loadNetworkFilters(); nf != nil {
		filter.Networks = nf
		hasAnyFilter = true
	}

	// Images
	if imf := loadImageFilters(); imf != nil {
		filter.Images = imf
		hasAnyFilter = true
	}

	if !hasAnyFilter {
		return nil
	}

	return filter
}

// loadVolumeFilters charge les filtres de volumes depuis l'environnement
func loadVolumeFilters() *filters.VolumeFilter {
	vf := &filters.VolumeFilter{}
	hasFilter := false

	if allowedNames := getEnvArray("VOLUMES__ALLOWED_NAMES"); len(allowedNames) > 0 {
		vf.AllowedNames = allowedNames
		hasFilter = true
	}

	if deniedNames := getEnvArray("VOLUMES__DENIED_NAMES"); len(deniedNames) > 0 {
		vf.DeniedNames = deniedNames
		hasFilter = true
	}

	if allowedPaths := getEnvArray("VOLUMES__ALLOWED_PATHS"); len(allowedPaths) > 0 {
		vf.AllowedPaths = allowedPaths
		hasFilter = true
	}

	if deniedPaths := getEnvArray("VOLUMES__DENIED_PATHS"); len(deniedPaths) > 0 {
		vf.DeniedPaths = deniedPaths
		hasFilter = true
	}

	if allowedDrivers := getEnvArray("VOLUMES__ALLOWED_DRIVERS"); len(allowedDrivers) > 0 {
		vf.AllowedDrivers = allowedDrivers
		hasFilter = true
	}

	if !hasFilter {
		return nil
	}
	return vf
}

// loadContainerFilters loads container filters from environment
func loadContainerFilters() *filters.ContainerFilter {
	cf := &filters.ContainerFilter{}
	hasFilter := false

	if allowedImages := getEnvArray("CONTAINERS__ALLOWED_IMAGES"); len(allowedImages) > 0 {
		cf.AllowedImages = allowedImages
		hasFilter = true
	}

	if deniedImages := getEnvArray("CONTAINERS__DENIED_IMAGES"); len(deniedImages) > 0 {
		cf.DeniedImages = deniedImages
		hasFilter = true
	}

	if allowedNames := getEnvArray("CONTAINERS__ALLOWED_NAMES"); len(allowedNames) > 0 {
		cf.AllowedNames = allowedNames
		hasFilter = true
	}

	if deniedNames := getEnvArray("CONTAINERS__DENIED_NAMES"); len(deniedNames) > 0 {
		cf.DeniedNames = deniedNames
		hasFilter = true
	}

	if requireLabels := getEnvMap("CONTAINERS__REQUIRE_LABELS"); len(requireLabels) > 0 {
		cf.RequireLabels = requireLabels
		hasFilter = true
	}

	if val := os.Getenv(envPrefix + "CONTAINERS__DENY_PRIVILEGED"); val != "" {
		cf.DenyPrivileged = parseBool(val)
		hasFilter = true
	}

	if val := os.Getenv(envPrefix + "CONTAINERS__DENY_HOST_NETWORK"); val != "" {
		cf.DenyHostNetwork = parseBool(val)
		hasFilter = true
	}

	if !hasFilter {
		return nil
	}
	return cf
}

// loadNetworkFilters charge les filtres de réseaux depuis l'environnement
func loadNetworkFilters() *filters.NetworkFilter {
	nf := &filters.NetworkFilter{}
	hasFilter := false

	if allowedNames := getEnvArray("NETWORKS__ALLOWED_NAMES"); len(allowedNames) > 0 {
		nf.AllowedNames = allowedNames
		hasFilter = true
	}

	if deniedNames := getEnvArray("NETWORKS__DENIED_NAMES"); len(deniedNames) > 0 {
		nf.DeniedNames = deniedNames
		hasFilter = true
	}

	if allowedDrivers := getEnvArray("NETWORKS__ALLOWED_DRIVERS"); len(allowedDrivers) > 0 {
		nf.AllowedDrivers = allowedDrivers
		hasFilter = true
	}

	if !hasFilter {
		return nil
	}
	return nf
}

// loadImageFilters charge les filtres d'images depuis l'environnement
func loadImageFilters() *filters.ImageFilter {
	imf := &filters.ImageFilter{}
	hasFilter := false

	if allowedRepos := getEnvArray("IMAGES__ALLOWED_REPOS"); len(allowedRepos) > 0 {
		imf.AllowedRepos = allowedRepos
		hasFilter = true
	}

	if deniedRepos := getEnvArray("IMAGES__DENIED_REPOS"); len(deniedRepos) > 0 {
		imf.DeniedRepos = deniedRepos
		hasFilter = true
	}

	if allowedTags := getEnvArray("IMAGES__ALLOWED_TAGS"); len(allowedTags) > 0 {
		imf.AllowedTags = allowedTags
		hasFilter = true
	}

	if deniedTags := getEnvArray("IMAGES__DENIED_TAGS"); len(deniedTags) > 0 {
		imf.DeniedTags = deniedTags
		hasFilter = true
	}

	if !hasFilter {
		return nil
	}
	return imf
}

// getEnvArray récupère une variable d'environnement et la convertit en tableau
// Format: "value1,value2,value3" ou "value1|value2|value3"
func getEnvArray(key string) []string {
	value := os.Getenv(envPrefix + key)
	if value == "" {
		return nil
	}

	// Support des séparateurs: virgule, pipe, point-virgule
	var items []string
	if strings.Contains(value, "|") {
		items = strings.Split(value, "|")
	} else if strings.Contains(value, ";") {
		items = strings.Split(value, ";")
	} else {
		items = strings.Split(value, ",")
	}

	// Trim spaces
	result := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// getEnvMap récupère une variable d'environnement et la convertit en map
// Format: "key1=value1,key2=value2" ou "key1=value1|key2=value2"
func getEnvMap(key string) map[string]string {
	value := os.Getenv(envPrefix + key)
	if value == "" {
		return nil
	}

	result := make(map[string]string)

	// Support des séparateurs: virgule, pipe, point-virgule
	var pairs []string
	if strings.Contains(value, "|") {
		pairs = strings.Split(value, "|")
	} else if strings.Contains(value, ";") {
		pairs = strings.Split(value, ";")
	} else {
		pairs = strings.Split(value, ",")
	}

	for _, pair := range pairs {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			if k != "" {
				result[k] = v
			}
		}
	}

	return result
}

// parseBool convertit une chaîne en booléen
func parseBool(value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	return value == "true" || value == "1" || value == "yes" || value == "on"
}

// MergeFilters fusionne les filtres en donnant priorité aux variables d'environnement
func MergeFilters(jsonFilter, envFilter *filters.AdvancedFilter) *filters.AdvancedFilter {
	// Si pas de filtre JSON, retourner le filtre env
	if jsonFilter == nil {
		return envFilter
	}

	// Si pas de filtre env, retourner le filtre JSON
	if envFilter == nil {
		return jsonFilter
	}

	// Fusionner avec priorité aux env vars
	result := &filters.AdvancedFilter{}

	// Volumes: env prioritaire
	if envFilter.Volumes != nil {
		result.Volumes = envFilter.Volumes
	} else {
		result.Volumes = jsonFilter.Volumes
	}

	// Containers: env prioritaire
	if envFilter.Containers != nil {
		result.Containers = envFilter.Containers
	} else {
		result.Containers = jsonFilter.Containers
	}

	// Networks: env prioritaire
	if envFilter.Networks != nil {
		result.Networks = envFilter.Networks
	} else {
		result.Networks = jsonFilter.Networks
	}

	// Images: env prioritaire
	if envFilter.Images != nil {
		result.Images = envFilter.Images
	} else {
		result.Images = jsonFilter.Images
	}

	return result
}
