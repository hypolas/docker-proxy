package config

import (
	"os"

	"docker-proxy/pkg/filters"
)

// GetDefaultFilters retourne les filtres par défaut pour la sécurité
func GetDefaultFilters() *filters.AdvancedFilter {
	// Obtenir le nom du conteneur docker-proxy depuis l'environnement
	proxyContainerName := getEnv("PROXY_CONTAINER_NAME", "docker-proxy")
	proxyNetworkName := getEnv("PROXY_NETWORK_NAME", "")

	containerFilter := &filters.ContainerFilter{
		// Interdire la manipulation du conteneur docker-proxy lui-même
		DeniedNames: []string{
			`^` + proxyContainerName + `$`,
			`^/` + proxyContainerName + `$`,
		},
	}

	volumeFilter := &filters.VolumeFilter{
		// Par défaut, interdire le montage du socket Docker
		DeniedPaths: []string{
			`^/var/run/docker\.sock$`,
			`^/run/docker\.sock$`,
		},
	}

	networkFilter := &filters.NetworkFilter{}
	if proxyNetworkName != "" {
		// Interdire la manipulation du réseau du proxy
		networkFilter.DeniedNames = []string{
			`^` + proxyNetworkName + `$`,
		}
	}

	return &filters.AdvancedFilter{
		Volumes:    volumeFilter,
		Containers: containerFilter,
		Networks:   networkFilter,
	}
}

// ApplyDefaults applique les filtres par défaut si aucun filtre n'est configuré
func ApplyDefaults(filter *filters.AdvancedFilter) *filters.AdvancedFilter {
	defaults := GetDefaultFilters()

	// Si aucun filtre n'est configuré, utiliser les défauts
	if filter == nil {
		return defaults
	}

	// Appliquer les défauts seulement pour les sections non configurées
	if filter.Volumes == nil {
		filter.Volumes = defaults.Volumes
	}

	if filter.Containers == nil {
		filter.Containers = defaults.Containers
	}

	if filter.Networks == nil && defaults.Networks != nil && len(defaults.Networks.DeniedNames) > 0 {
		filter.Networks = defaults.Networks
	}

	return filter
}

// CanOverrideDefaults vérifie si l'utilisateur peut désactiver les défauts
func CanOverrideDefaults() bool {
	// Permettre de désactiver les défauts via variable d'environnement
	val := os.Getenv("DKRPRX__DISABLE_DEFAULTS")
	if val == "" {
		return false
	}
	return parseBool(val)
}
