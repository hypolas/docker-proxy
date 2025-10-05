package filters

import (
	"encoding/json"
	"testing"
)

func TestCheckVolumeMount(t *testing.T) {
	tests := []struct {
		name          string
		filter        *AdvancedFilter
		volumeName    string
		hostPath      string
		driver        string
		expectAllowed bool
		expectReason  string
	}{
		{
			name:          "No filter returns allowed",
			filter:        &AdvancedFilter{},
			volumeName:    "test-volume",
			hostPath:      "/data",
			driver:        "local",
			expectAllowed: true,
			expectReason:  "",
		},
		{
			name: "Denied path blocked",
			filter: &AdvancedFilter{
				Volumes: &VolumeFilter{
					DeniedPaths: []string{`^/var/run/docker\.sock$`},
				},
			},
			volumeName:    "socket",
			hostPath:      "/var/run/docker.sock",
			driver:        "local",
			expectAllowed: false,
			expectReason:  "host path is denied: /var/run/docker.sock",
		},
		{
			name: "Allowed path passes",
			filter: &AdvancedFilter{
				Volumes: &VolumeFilter{
					AllowedPaths: []string{`^/data/.*`},
				},
			},
			volumeName:    "data",
			hostPath:      "/data/app",
			driver:        "local",
			expectAllowed: true,
			expectReason:  "",
		},
		{
			name: "Path not in allowed list blocked",
			filter: &AdvancedFilter{
				Volumes: &VolumeFilter{
					AllowedPaths: []string{`^/data/.*`},
				},
			},
			volumeName:    "etc",
			hostPath:      "/etc/passwd",
			driver:        "local",
			expectAllowed: false,
			expectReason:  "host path not in allowed list: /etc/passwd",
		},
		{
			name: "Denied volume name blocked",
			filter: &AdvancedFilter{
				Volumes: &VolumeFilter{
					DeniedNames: []string{`^prod-.*`},
				},
			},
			volumeName:    "prod-db",
			hostPath:      "/data",
			driver:        "local",
			expectAllowed: false,
			expectReason:  "volume name is denied: prod-db",
		},
		{
			name: "Allowed volume name passes",
			filter: &AdvancedFilter{
				Volumes: &VolumeFilter{
					AllowedNames: []string{`^app-.*`},
				},
			},
			volumeName:    "app-data",
			hostPath:      "/data",
			driver:        "local",
			expectAllowed: true,
			expectReason:  "",
		},
		{
			name: "Volume name not in allowed list blocked",
			filter: &AdvancedFilter{
				Volumes: &VolumeFilter{
					AllowedNames: []string{`^app-.*`},
				},
			},
			volumeName:    "db-data",
			hostPath:      "/data",
			driver:        "local",
			expectAllowed: false,
			expectReason:  "volume name not in allowed list: db-data",
		},
		{
			name: "Driver not allowed blocked",
			filter: &AdvancedFilter{
				Volumes: &VolumeFilter{
					AllowedDrivers: []string{"local"},
				},
			},
			volumeName:    "nfs-volume",
			hostPath:      "",
			driver:        "nfs",
			expectAllowed: false,
			expectReason:  "volume driver not allowed: nfs",
		},
		{
			name: "Allowed driver passes",
			filter: &AdvancedFilter{
				Volumes: &VolumeFilter{
					AllowedDrivers: []string{"local", "nfs"},
				},
			},
			volumeName:    "nfs-volume",
			hostPath:      "",
			driver:        "nfs",
			expectAllowed: true,
			expectReason:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, reason := tt.filter.CheckVolumeMount(tt.volumeName, tt.hostPath, tt.driver)
			if allowed != tt.expectAllowed {
				t.Errorf("Expected allowed=%v, got %v", tt.expectAllowed, allowed)
			}
			if reason != tt.expectReason {
				t.Errorf("Expected reason='%s', got '%s'", tt.expectReason, reason)
			}
		})
	}
}

func TestCheckContainerCreate(t *testing.T) {
	tests := []struct {
		name          string
		filter        *AdvancedFilter
		image         string
		containerName string
		config        map[string]interface{}
		expectAllowed bool
		expectReason  string
	}{
		{
			name:          "No filter returns allowed",
			filter:        &AdvancedFilter{},
			image:         "nginx:latest",
			containerName: "web",
			config:        map[string]interface{}{},
			expectAllowed: true,
			expectReason:  "",
		},
		{
			name: "Denied image blocked",
			filter: &AdvancedFilter{
				Containers: &ContainerFilter{
					DeniedImages: []string{`.*:latest$`},
				},
			},
			image:         "nginx:latest",
			containerName: "web",
			config:        map[string]interface{}{},
			expectAllowed: false,
			expectReason:  "image is denied: nginx:latest",
		},
		{
			name: "Allowed image passes",
			filter: &AdvancedFilter{
				Containers: &ContainerFilter{
					AllowedImages: []string{`^registry\.company\.com/.*`},
				},
			},
			image:         "registry.company.com/app:v1.0",
			containerName: "app",
			config:        map[string]interface{}{},
			expectAllowed: true,
			expectReason:  "",
		},
		{
			name: "Image not in allowed list blocked",
			filter: &AdvancedFilter{
				Containers: &ContainerFilter{
					AllowedImages: []string{`^registry\.company\.com/.*`},
				},
			},
			image:         "docker.io/nginx:latest",
			containerName: "web",
			config:        map[string]interface{}{},
			expectAllowed: false,
			expectReason:  "image not in allowed list: docker.io/nginx:latest",
		},
		{
			name: "Privileged container denied",
			filter: &AdvancedFilter{
				Containers: &ContainerFilter{
					DenyPrivileged: true,
				},
			},
			image:         "nginx:latest",
			containerName: "web",
			config: map[string]interface{}{
				"HostConfig": map[string]interface{}{
					"Privileged": true,
				},
			},
			expectAllowed: false,
			expectReason:  "privileged containers are denied",
		},
		{
			name: "Non-privileged container allowed",
			filter: &AdvancedFilter{
				Containers: &ContainerFilter{
					DenyPrivileged: true,
				},
			},
			image:         "nginx:latest",
			containerName: "web",
			config: map[string]interface{}{
				"HostConfig": map[string]interface{}{
					"Privileged": false,
				},
			},
			expectAllowed: true,
			expectReason:  "",
		},
		{
			name: "Host network denied",
			filter: &AdvancedFilter{
				Containers: &ContainerFilter{
					DenyHostNetwork: true,
				},
			},
			image:         "nginx:latest",
			containerName: "web",
			config: map[string]interface{}{
				"HostConfig": map[string]interface{}{
					"NetworkMode": "host",
				},
			},
			expectAllowed: false,
			expectReason:  "host network mode is denied",
		},
		{
			name: "Bridge network allowed",
			filter: &AdvancedFilter{
				Containers: &ContainerFilter{
					DenyHostNetwork: true,
				},
			},
			image:         "nginx:latest",
			containerName: "web",
			config: map[string]interface{}{
				"HostConfig": map[string]interface{}{
					"NetworkMode": "bridge",
				},
			},
			expectAllowed: true,
			expectReason:  "",
		},
		{
			name: "Required labels missing",
			filter: &AdvancedFilter{
				Containers: &ContainerFilter{
					RequireLabels: map[string]string{
						"env":  "production",
						"team": "backend",
					},
				},
			},
			image:         "nginx:latest",
			containerName: "web",
			config:        map[string]interface{}{},
			expectAllowed: false,
			expectReason:  "required labels are missing",
		},
		{
			name: "Required labels present",
			filter: &AdvancedFilter{
				Containers: &ContainerFilter{
					RequireLabels: map[string]string{
						"env":  "production",
						"team": "backend",
					},
				},
			},
			image:         "nginx:latest",
			containerName: "web",
			config: map[string]interface{}{
				"Labels": map[string]interface{}{
					"env":  "production",
					"team": "backend",
				},
			},
			expectAllowed: true,
			expectReason:  "",
		},
		{
			name: "Required label mismatch",
			filter: &AdvancedFilter{
				Containers: &ContainerFilter{
					RequireLabels: map[string]string{
						"env": "production",
					},
				},
			},
			image:         "nginx:latest",
			containerName: "web",
			config: map[string]interface{}{
				"Labels": map[string]interface{}{
					"env": "development",
				},
			},
			expectAllowed: false,
			expectReason:  "required label missing or mismatch: env",
		},
		{
			name: "Container name denied",
			filter: &AdvancedFilter{
				Containers: &ContainerFilter{
					DeniedNames: []string{`^prod-.*`},
				},
			},
			image:         "nginx:latest",
			containerName: "prod-db",
			config:        map[string]interface{}{},
			expectAllowed: false,
			expectReason:  "container name is denied: prod-db",
		},
		{
			name: "Container name allowed",
			filter: &AdvancedFilter{
				Containers: &ContainerFilter{
					AllowedNames: []string{`^app-.*`},
				},
			},
			image:         "nginx:latest",
			containerName: "app-web",
			config:        map[string]interface{}{},
			expectAllowed: true,
			expectReason:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, reason := tt.filter.CheckContainerCreate(tt.image, tt.containerName, tt.config)
			if allowed != tt.expectAllowed {
				t.Errorf("Expected allowed=%v, got %v", tt.expectAllowed, allowed)
			}
			if reason != tt.expectReason {
				t.Errorf("Expected reason='%s', got '%s'", tt.expectReason, reason)
			}
		})
	}
}

func TestCheckNetworkCreate(t *testing.T) {
	tests := []struct {
		name          string
		filter        *AdvancedFilter
		networkName   string
		driver        string
		expectAllowed bool
		expectReason  string
	}{
		{
			name:          "No filter returns allowed",
			filter:        &AdvancedFilter{},
			networkName:   "app-network",
			driver:        "bridge",
			expectAllowed: true,
			expectReason:  "",
		},
		{
			name: "Denied network name blocked",
			filter: &AdvancedFilter{
				Networks: &NetworkFilter{
					DeniedNames: []string{`^host$`},
				},
			},
			networkName:   "host",
			driver:        "host",
			expectAllowed: false,
			expectReason:  "network name is denied: host",
		},
		{
			name: "Allowed network name passes",
			filter: &AdvancedFilter{
				Networks: &NetworkFilter{
					AllowedNames: []string{`^app-.*`},
				},
			},
			networkName:   "app-backend",
			driver:        "bridge",
			expectAllowed: true,
			expectReason:  "",
		},
		{
			name: "Network name not in allowed list blocked",
			filter: &AdvancedFilter{
				Networks: &NetworkFilter{
					AllowedNames: []string{`^app-.*`},
				},
			},
			networkName:   "db-network",
			driver:        "bridge",
			expectAllowed: false,
			expectReason:  "network name not in allowed list: db-network",
		},
		{
			name: "Driver not allowed blocked",
			filter: &AdvancedFilter{
				Networks: &NetworkFilter{
					AllowedDrivers: []string{"bridge", "overlay"},
				},
			},
			networkName:   "macvlan-net",
			driver:        "macvlan",
			expectAllowed: false,
			expectReason:  "network driver not allowed: macvlan",
		},
		{
			name: "Allowed driver passes",
			filter: &AdvancedFilter{
				Networks: &NetworkFilter{
					AllowedDrivers: []string{"bridge", "overlay"},
				},
			},
			networkName:   "overlay-net",
			driver:        "overlay",
			expectAllowed: true,
			expectReason:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, reason := tt.filter.CheckNetworkCreate(tt.networkName, tt.driver)
			if allowed != tt.expectAllowed {
				t.Errorf("Expected allowed=%v, got %v", tt.expectAllowed, allowed)
			}
			if reason != tt.expectReason {
				t.Errorf("Expected reason='%s', got '%s'", tt.expectReason, reason)
			}
		})
	}
}

func TestCheckImageOperation(t *testing.T) {
	tests := []struct {
		name          string
		filter        *AdvancedFilter
		imageName     string
		expectAllowed bool
		expectReason  string
	}{
		{
			name:          "No filter returns allowed",
			filter:        &AdvancedFilter{},
			imageName:     "nginx:latest",
			expectAllowed: true,
			expectReason:  "",
		},
		{
			name: "Denied repository blocked",
			filter: &AdvancedFilter{
				Images: &ImageFilter{
					DeniedRepos: []string{`.*\.suspicious\.com/.*`},
				},
			},
			imageName:     "registry.suspicious.com/app:v1.0",
			expectAllowed: false,
			expectReason:  "image repository is denied: registry.suspicious.com/app",
		},
		{
			name: "Allowed repository passes",
			filter: &AdvancedFilter{
				Images: &ImageFilter{
					AllowedRepos: []string{`^registry\.company\.com/.*`},
				},
			},
			imageName:     "registry.company.com/app:v1.0",
			expectAllowed: true,
			expectReason:  "",
		},
		{
			name: "Repository not in allowed list blocked",
			filter: &AdvancedFilter{
				Images: &ImageFilter{
					AllowedRepos: []string{`^registry\.company\.com/.*`},
				},
			},
			imageName:     "docker.io/nginx:latest",
			expectAllowed: false,
			expectReason:  "image repository not in allowed list: docker.io/nginx",
		},
		{
			name: "Denied tag blocked",
			filter: &AdvancedFilter{
				Images: &ImageFilter{
					DeniedTags: []string{`^latest$`},
				},
			},
			imageName:     "nginx:latest",
			expectAllowed: false,
			expectReason:  "image tag is denied: latest",
		},
		{
			name: "Allowed tag passes",
			filter: &AdvancedFilter{
				Images: &ImageFilter{
					AllowedTags: []string{`^v[0-9]+\.[0-9]+\.[0-9]+$`},
				},
			},
			imageName:     "nginx:v1.2.3",
			expectAllowed: true,
			expectReason:  "",
		},
		{
			name: "Tag not in allowed list blocked",
			filter: &AdvancedFilter{
				Images: &ImageFilter{
					AllowedTags: []string{`^v[0-9]+\.[0-9]+\.[0-9]+$`},
				},
			},
			imageName:     "nginx:latest",
			expectAllowed: false,
			expectReason:  "image tag not in allowed list: latest",
		},
		{
			name: "Image without tag defaults to latest",
			filter: &AdvancedFilter{
				Images: &ImageFilter{
					DeniedTags: []string{`^latest$`},
				},
			},
			imageName:     "nginx",
			expectAllowed: false,
			expectReason:  "image tag is denied: latest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, reason := tt.filter.CheckImageOperation(tt.imageName)
			if allowed != tt.expectAllowed {
				t.Errorf("Expected allowed=%v, got %v", tt.expectAllowed, allowed)
			}
			if reason != tt.expectReason {
				t.Errorf("Expected reason='%s', got '%s'", tt.expectReason, reason)
			}
		})
	}
}

func TestParseImageName(t *testing.T) {
	tests := []struct {
		name        string
		imageName   string
		expectedRepo string
		expectedTag  string
	}{
		{
			name:         "Image with tag",
			imageName:    "nginx:1.21",
			expectedRepo: "nginx",
			expectedTag:  "1.21",
		},
		{
			name:         "Image without tag",
			imageName:    "nginx",
			expectedRepo: "nginx",
			expectedTag:  "latest",
		},
		{
			name:         "Registry with tag",
			imageName:    "registry.company.com/app:v1.0",
			expectedRepo: "registry.company.com/app",
			expectedTag:  "v1.0",
		},
		{
			name:         "Registry without tag",
			imageName:    "registry.company.com/app",
			expectedRepo: "registry.company.com/app",
			expectedTag:  "latest",
		},
		{
			name:         "Docker Hub with namespace",
			imageName:    "library/nginx:alpine",
			expectedRepo: "library/nginx",
			expectedTag:  "alpine",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, tag := parseImageName(tt.imageName)
			if repo != tt.expectedRepo {
				t.Errorf("Expected repo '%s', got '%s'", tt.expectedRepo, repo)
			}
			if tag != tt.expectedTag {
				t.Errorf("Expected tag '%s', got '%s'", tt.expectedTag, tag)
			}
		})
	}
}

func TestLoadFromJSON(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		expectError bool
	}{
		{
			name: "Valid JSON",
			jsonData: `{
				"volumes": {
					"denied_paths": ["/var/run/docker.sock"]
				},
				"containers": {
					"deny_privileged": true
				}
			}`,
			expectError: false,
		},
		{
			name:        "Invalid JSON",
			jsonData:    `{invalid json}`,
			expectError: true,
		},
		{
			name:        "Empty JSON",
			jsonData:    `{}`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := LoadFromJSON([]byte(tt.jsonData))
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if filter == nil {
					t.Error("Expected non-nil filter")
				}
			}
		})
	}
}

func TestAdvancedFilterJSON(t *testing.T) {
	filter := &AdvancedFilter{
		Volumes: &VolumeFilter{
			DeniedPaths: []string{"/var/run/docker.sock", "/etc"},
		},
		Containers: &ContainerFilter{
			DenyPrivileged:  true,
			DenyHostNetwork: true,
			AllowedImages:   []string{"^registry\\.company\\.com/.*"},
		},
		Networks: &NetworkFilter{
			AllowedDrivers: []string{"bridge", "overlay"},
		},
		Images: &ImageFilter{
			DeniedTags: []string{"^latest$"},
		},
	}

	// Test marshaling
	data, err := json.Marshal(filter)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Test unmarshaling
	var unmarshaled AdvancedFilter
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify some fields
	if unmarshaled.Containers.DenyPrivileged != true {
		t.Error("DenyPrivileged not preserved")
	}
	if len(unmarshaled.Volumes.DeniedPaths) != 2 {
		t.Errorf("Expected 2 denied paths, got %d", len(unmarshaled.Volumes.DeniedPaths))
	}
}
