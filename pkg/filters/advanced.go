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
	AllowedNames   []string `json:"allowed_names,omitempty"`   // Liste blanche des noms
	DeniedNames    []string `json:"denied_names,omitempty"`    // Liste noire des noms
	AllowedPaths   []string `json:"allowed_paths,omitempty"`   // Chemins autorisés (patterns)
	DeniedPaths    []string `json:"denied_paths,omitempty"`    // Chemins interdits (patterns)
	AllowedDrivers []string `json:"allowed_drivers,omitempty"` // Drivers autorisés
}

// ContainerFilter defines filtering rules for containers
type ContainerFilter struct {
	AllowedImages   []string          `json:"allowed_images,omitempty"`    // Allowed images (patterns)
	DeniedImages    []string          `json:"denied_images,omitempty"`     // Denied images (patterns)
	AllowedNames    []string          `json:"allowed_names,omitempty"`     // Allowed names (patterns)
	DeniedNames     []string          `json:"denied_names,omitempty"`      // Denied names (patterns)
	RequireLabels   map[string]string `json:"require_labels,omitempty"`    // Required labels
	DenyPrivileged  bool              `json:"deny_privileged,omitempty"`   // Deny privileged containers
	DenyHostNetwork bool              `json:"deny_host_network,omitempty"` // Deny host network
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

// CheckVolumeMount checks if a volume mount is allowed
func (af *AdvancedFilter) CheckVolumeMount(volumeName, hostPath, driver string) (bool, string) {
	if af.Volumes == nil {
		return true, ""
	}

	vf := af.Volumes

	// Check denied names
	if ok, msg := checkDeniedList(vf.DeniedNames, volumeName, "volume name is denied"); !ok {
		return false, msg
	}

	// Check allowed names
	if ok, msg := checkAllowedList(vf.AllowedNames, volumeName, "volume name not in allowed list"); !ok {
		return false, msg
	}

	// Check denied paths
	if hostPath != "" {
		if ok, msg := checkDeniedList(vf.DeniedPaths, hostPath, "host path is denied"); !ok {
			return false, msg
		}

		// Check allowed paths
		if ok, msg := checkAllowedList(vf.AllowedPaths, hostPath, "host path not in allowed list"); !ok {
			return false, msg
		}
	}

	// Check driver
	if driver != "" && len(vf.AllowedDrivers) > 0 {
		if !contains(vf.AllowedDrivers, driver) {
			return false, "volume driver not allowed: " + driver
		}
	}

	return true, ""
}

// CheckContainerCreate checks if a container creation is allowed
func (af *AdvancedFilter) CheckContainerCreate(image, name string, config map[string]interface{}) (bool, string) {
	if af.Containers == nil {
		return true, ""
	}

	cf := af.Containers

	// Check denied images
	if ok, msg := checkDeniedList(cf.DeniedImages, image, "image is denied"); !ok {
		return false, msg
	}

	// Check allowed images
	if ok, msg := checkAllowedList(cf.AllowedImages, image, "image not in allowed list"); !ok {
		return false, msg
	}

	// Check container names
	if name != "" {
		if ok, msg := checkDeniedList(cf.DeniedNames, name, "container name is denied"); !ok {
			return false, msg
		}
		if ok, msg := checkAllowedList(cf.AllowedNames, name, "container name not in allowed list"); !ok {
			return false, msg
		}
	}

	// Check HostConfig settings
	if ok, msg := checkHostConfig(cf, config); !ok {
		return false, msg
	}

	// Check required labels
	if ok, msg := checkRequiredLabels(cf.RequireLabels, config); !ok {
		return false, msg
	}

	return true, ""
}

// checkHostConfig validates HostConfig settings (privileged, host network)
func checkHostConfig(cf *ContainerFilter, config map[string]interface{}) (bool, string) {
	hostConfig, ok := config["HostConfig"].(map[string]interface{})
	if !ok {
		return true, ""
	}

	// Check privileged mode
	if cf.DenyPrivileged {
		if privileged, ok := hostConfig["Privileged"].(bool); ok && privileged {
			return false, "privileged containers are denied"
		}
	}

	// Check host network
	if cf.DenyHostNetwork {
		if networkMode, ok := hostConfig["NetworkMode"].(string); ok && networkMode == "host" {
			return false, "host network mode is denied"
		}
	}

	return true, ""
}

// checkRequiredLabels validates that all required labels are present
func checkRequiredLabels(requiredLabels map[string]string, config map[string]interface{}) (bool, string) {
	if len(requiredLabels) == 0 {
		return true, ""
	}

	labels, ok := config["Labels"].(map[string]interface{})
	if !ok {
		return false, "required labels are missing"
	}

	for key, value := range requiredLabels {
		if labelValue, ok := labels[key].(string); !ok || labelValue != value {
			return false, "required label missing or mismatch: " + key
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

// CheckImageOperation checks if an image operation is allowed
func (af *AdvancedFilter) CheckImageOperation(imageName string) (bool, string) {
	if af.Images == nil {
		return true, ""
	}

	imf := af.Images
	repo, tag := parseImageName(imageName)

	// Check repository filters
	if ok, msg := checkDeniedList(imf.DeniedRepos, repo, "image repository is denied"); !ok {
		return false, msg
	}
	if ok, msg := checkAllowedList(imf.AllowedRepos, repo, "image repository not in allowed list"); !ok {
		return false, msg
	}

	// Check tag filters
	if tag != "" {
		if ok, msg := checkDeniedList(imf.DeniedTags, tag, "image tag is denied"); !ok {
			return false, msg
		}
		if ok, msg := checkAllowedList(imf.AllowedTags, tag, "image tag not in allowed list"); !ok {
			return false, msg
		}
	}

	return true, ""
}

// checkDeniedList checks if a value matches any pattern in the denied list
func checkDeniedList(deniedList []string, value, errorPrefix string) (bool, string) {
	if len(deniedList) == 0 {
		return true, ""
	}

	for _, denied := range deniedList {
		if matched, _ := regexp.MatchString(denied, value); matched {
			return false, errorPrefix + ": " + value
		}
	}
	return true, ""
}

// checkAllowedList checks if a value matches at least one pattern in the allowed list
func checkAllowedList(allowedList []string, value, errorPrefix string) (bool, string) {
	if len(allowedList) == 0 {
		return true, ""
	}

	for _, allow := range allowedList {
		if matched, _ := regexp.MatchString(allow, value); matched {
			return true, ""
		}
	}
	return false, errorPrefix + ": " + value
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
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
