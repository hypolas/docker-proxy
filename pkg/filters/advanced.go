package filters

import (
	"encoding/json"
	"regexp"
	"strings"
)

// AdvancedFilter définit des règles de filtrage avancées
type AdvancedFilter struct {
	Volumes    *VolumeFilter    `json:"volumes,omitempty"`
	Containers *ContainerFilter `json:"containers,omitempty"`
	Networks   *NetworkFilter   `json:"networks,omitempty"`
	Images     *ImageFilter     `json:"images,omitempty"`
}

// VolumeFilter définit les règles de filtrage pour les volumes
type VolumeFilter struct {
	AllowedNames  []string `json:"allowed_names,omitempty"`   // Liste blanche des noms
	DeniedNames   []string `json:"denied_names,omitempty"`    // Liste noire des noms
	AllowedPaths  []string `json:"allowed_paths,omitempty"`   // Chemins autorisés (patterns)
	DeniedPaths   []string `json:"denied_paths,omitempty"`    // Chemins interdits (patterns)
	AllowedDrivers []string `json:"allowed_drivers,omitempty"` // Drivers autorisés
}

// ContainerFilter définit les règles de filtrage pour les conteneurs
type ContainerFilter struct {
	AllowedImages  []string          `json:"allowed_images,omitempty"`  // Images autorisées (patterns)
	DeniedImages   []string          `json:"denied_images,omitempty"`   // Images interdites (patterns)
	AllowedNames   []string          `json:"allowed_names,omitempty"`   // Noms autorisés (patterns)
	DeniedNames    []string          `json:"denied_names,omitempty"`    // Noms interdits (patterns)
	RequireLabels  map[string]string `json:"require_labels,omitempty"`  // Labels requis
	DenyPrivileged bool              `json:"deny_privileged,omitempty"` // Interdire les conteneurs privilégiés
	DenyHostNetwork bool             `json:"deny_host_network,omitempty"` // Interdire host network
}

// NetworkFilter définit les règles de filtrage pour les réseaux
type NetworkFilter struct {
	AllowedNames   []string `json:"allowed_names,omitempty"`   // Noms autorisés (patterns)
	DeniedNames    []string `json:"denied_names,omitempty"`    // Noms interdits (patterns)
	AllowedDrivers []string `json:"allowed_drivers,omitempty"` // Drivers autorisés
}

// ImageFilter définit les règles de filtrage pour les images
type ImageFilter struct {
	AllowedRepos []string `json:"allowed_repos,omitempty"` // Registres/repos autorisés (patterns)
	DeniedRepos  []string `json:"denied_repos,omitempty"`  // Registres/repos interdits (patterns)
	AllowedTags  []string `json:"allowed_tags,omitempty"`  // Tags autorisés (patterns)
	DeniedTags   []string `json:"denied_tags,omitempty"`   // Tags interdits (patterns)
}

// CheckVolumeMount vérifie si un montage de volume est autorisé
func (af *AdvancedFilter) CheckVolumeMount(volumeName, hostPath, driver string) (bool, string) {
	if af.Volumes == nil {
		return true, ""
	}

	vf := af.Volumes

	// Vérifier la liste noire des noms
	if len(vf.DeniedNames) > 0 {
		for _, denied := range vf.DeniedNames {
			if matched, _ := regexp.MatchString(denied, volumeName); matched {
				return false, "volume name is denied: " + volumeName
			}
		}
	}

	// Vérifier la liste blanche des noms (si définie)
	if len(vf.AllowedNames) > 0 {
		allowed := false
		for _, allow := range vf.AllowedNames {
			if matched, _ := regexp.MatchString(allow, volumeName); matched {
				allowed = true
				break
			}
		}
		if !allowed {
			return false, "volume name not in allowed list: " + volumeName
		}
	}

	// Vérifier les chemins interdits
	if len(vf.DeniedPaths) > 0 && hostPath != "" {
		for _, denied := range vf.DeniedPaths {
			if matched, _ := regexp.MatchString(denied, hostPath); matched {
				return false, "host path is denied: " + hostPath
			}
		}
	}

	// Vérifier les chemins autorisés
	if len(vf.AllowedPaths) > 0 && hostPath != "" {
		allowed := false
		for _, allow := range vf.AllowedPaths {
			if matched, _ := regexp.MatchString(allow, hostPath); matched {
				allowed = true
				break
			}
		}
		if !allowed {
			return false, "host path not in allowed list: " + hostPath
		}
	}

	// Vérifier le driver
	if len(vf.AllowedDrivers) > 0 && driver != "" {
		allowed := false
		for _, allowedDriver := range vf.AllowedDrivers {
			if allowedDriver == driver {
				allowed = true
				break
			}
		}
		if !allowed {
			return false, "volume driver not allowed: " + driver
		}
	}

	return true, ""
}

// CheckContainerCreate vérifie si la création d'un conteneur est autorisée
func (af *AdvancedFilter) CheckContainerCreate(image, name string, config map[string]interface{}) (bool, string) {
	if af.Containers == nil {
		return true, ""
	}

	cf := af.Containers

	// Vérifier les images interdites
	if len(cf.DeniedImages) > 0 {
		for _, denied := range cf.DeniedImages {
			if matched, _ := regexp.MatchString(denied, image); matched {
				return false, "image is denied: " + image
			}
		}
	}

	// Vérifier les images autorisées
	if len(cf.AllowedImages) > 0 {
		allowed := false
		for _, allow := range cf.AllowedImages {
			if matched, _ := regexp.MatchString(allow, image); matched {
				allowed = true
				break
			}
		}
		if !allowed {
			return false, "image not in allowed list: " + image
		}
	}

	// Vérifier les noms interdits
	if len(cf.DeniedNames) > 0 && name != "" {
		for _, denied := range cf.DeniedNames {
			if matched, _ := regexp.MatchString(denied, name); matched {
				return false, "container name is denied: " + name
			}
		}
	}

	// Vérifier les noms autorisés
	if len(cf.AllowedNames) > 0 && name != "" {
		allowed := false
		for _, allow := range cf.AllowedNames {
			if matched, _ := regexp.MatchString(allow, name); matched {
				allowed = true
				break
			}
		}
		if !allowed {
			return false, "container name not in allowed list: " + name
		}
	}

	// Vérifier mode privilégié
	if cf.DenyPrivileged {
		if hostConfig, ok := config["HostConfig"].(map[string]interface{}); ok {
			if privileged, ok := hostConfig["Privileged"].(bool); ok && privileged {
				return false, "privileged containers are denied"
			}
		}
	}

	// Vérifier host network
	if cf.DenyHostNetwork {
		if hostConfig, ok := config["HostConfig"].(map[string]interface{}); ok {
			if networkMode, ok := hostConfig["NetworkMode"].(string); ok && networkMode == "host" {
				return false, "host network mode is denied"
			}
		}
	}

	// Vérifier les labels requis
	if len(cf.RequireLabels) > 0 {
		labels, ok := config["Labels"].(map[string]interface{})
		if !ok {
			return false, "required labels are missing"
		}
		for key, value := range cf.RequireLabels {
			if labelValue, ok := labels[key].(string); !ok || labelValue != value {
				return false, "required label missing or mismatch: " + key
			}
		}
	}

	return true, ""
}

// CheckNetworkCreate vérifie si la création d'un réseau est autorisée
func (af *AdvancedFilter) CheckNetworkCreate(name, driver string) (bool, string) {
	if af.Networks == nil {
		return true, ""
	}

	nf := af.Networks

	// Vérifier les noms interdits
	if len(nf.DeniedNames) > 0 {
		for _, denied := range nf.DeniedNames {
			if matched, _ := regexp.MatchString(denied, name); matched {
				return false, "network name is denied: " + name
			}
		}
	}

	// Vérifier les noms autorisés
	if len(nf.AllowedNames) > 0 {
		allowed := false
		for _, allow := range nf.AllowedNames {
			if matched, _ := regexp.MatchString(allow, name); matched {
				allowed = true
				break
			}
		}
		if !allowed {
			return false, "network name not in allowed list: " + name
		}
	}

	// Vérifier le driver
	if len(nf.AllowedDrivers) > 0 && driver != "" {
		allowed := false
		for _, allowedDriver := range nf.AllowedDrivers {
			if allowedDriver == driver {
				allowed = true
				break
			}
		}
		if !allowed {
			return false, "network driver not allowed: " + driver
		}
	}

	return true, ""
}

// CheckImageOperation vérifie si une opération sur une image est autorisée
func (af *AdvancedFilter) CheckImageOperation(imageName string) (bool, string) {
	if af.Images == nil {
		return true, ""
	}

	imf := af.Images

	// Extraire repo et tag
	repo, tag := parseImageName(imageName)

	// Vérifier les repos interdits
	if len(imf.DeniedRepos) > 0 {
		for _, denied := range imf.DeniedRepos {
			if matched, _ := regexp.MatchString(denied, repo); matched {
				return false, "image repository is denied: " + repo
			}
		}
	}

	// Vérifier les repos autorisés
	if len(imf.AllowedRepos) > 0 {
		allowed := false
		for _, allow := range imf.AllowedRepos {
			if matched, _ := regexp.MatchString(allow, repo); matched {
				allowed = true
				break
			}
		}
		if !allowed {
			return false, "image repository not in allowed list: " + repo
		}
	}

	// Vérifier les tags interdits
	if len(imf.DeniedTags) > 0 && tag != "" {
		for _, denied := range imf.DeniedTags {
			if matched, _ := regexp.MatchString(denied, tag); matched {
				return false, "image tag is denied: " + tag
			}
		}
	}

	// Vérifier les tags autorisés
	if len(imf.AllowedTags) > 0 && tag != "" {
		allowed := false
		for _, allow := range imf.AllowedTags {
			if matched, _ := regexp.MatchString(allow, tag); matched {
				allowed = true
				break
			}
		}
		if !allowed {
			return false, "image tag not in allowed list: " + tag
		}
	}

	return true, ""
}

// parseImageName extrait le repository et le tag d'un nom d'image
func parseImageName(imageName string) (repo, tag string) {
	parts := strings.Split(imageName, ":")
	repo = parts[0]
	tag = "latest"
	if len(parts) > 1 {
		tag = parts[1]
	}
	return
}

// LoadFromJSON charge les filtres depuis un JSON
func LoadFromJSON(jsonData []byte) (*AdvancedFilter, error) {
	var filter AdvancedFilter
	if err := json.Unmarshal(jsonData, &filter); err != nil {
		return nil, err
	}
	return &filter, nil
}
