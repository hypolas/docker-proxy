package config

import (
	"context"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

// DetectDockerAPIVersion détecte automatiquement la version de l'API Docker
func DetectDockerAPIVersion(dockerSocket string, logger *logrus.Logger) string {
	// Créer un client Docker
	opts := []client.Opt{
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	}

	// Si un socket spécifique est fourni, l'utiliser
	if dockerSocket != "" && strings.HasPrefix(dockerSocket, "unix://") {
		socketPath := strings.TrimPrefix(dockerSocket, "unix://")
		opts = append(opts, client.WithHost("unix://"+socketPath))
	} else if dockerSocket != "" {
		opts = append(opts, client.WithHost(dockerSocket))
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		logger.Warnf("Failed to create Docker client for version detection: %v", err)
		return "v1.41" // Fallback version
	}
	defer cli.Close()

	// Timeout pour la détection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Récupérer la version du serveur Docker
	version, err := cli.ServerVersion(ctx)
	if err != nil {
		logger.Warnf("Failed to detect Docker API version: %v", err)
		return "v1.41" // Fallback version
	}

	apiVersion := version.APIVersion
	logger.Infof("Detected Docker API version: %s", apiVersion)

	return apiVersion
}
