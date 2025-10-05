package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"

	"docker-proxy/pkg/filters"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AdvancedFilterMiddleware crée un middleware pour les filtres avancés
func AdvancedFilterMiddleware(filter *filters.AdvancedFilter, logger *logrus.Logger) gin.HandlerFunc {
	if filter == nil {
		// Pas de filtres avancés configurés
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.Request.URL.Path

		// Filtrer uniquement les opérations de création/modification
		if method != "POST" && method != "PUT" {
			c.Next()
			return
		}

		// Déterminer le type d'opération
		if matched, _ := regexp.MatchString(`/containers/create`, path); matched {
			if !checkContainerCreate(c, filter, logger) {
				return
			}
		} else if matched, _ := regexp.MatchString(`/volumes/create`, path); matched {
			if !checkVolumeCreate(c, filter, logger) {
				return
			}
		} else if matched, _ := regexp.MatchString(`/networks/create`, path); matched {
			if !checkNetworkCreate(c, filter, logger) {
				return
			}
		} else if matched, _ := regexp.MatchString(`/images/create`, path); matched {
			if !checkImageCreate(c, filter, logger) {
				return
			}
		} else if matched, _ := regexp.MatchString(`/build`, path); matched {
			// Vérifier le tag de l'image dans les query params
			if !checkImageBuild(c, filter, logger) {
				return
			}
		}

		c.Next()
	}
}

// checkContainerCreate vérifie la création de conteneur
func checkContainerCreate(c *gin.Context, filter *filters.AdvancedFilter, logger *logrus.Logger) bool {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return false
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var config map[string]interface{}
	if err := json.Unmarshal(body, &config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return false
	}

	// Extraire l'image et le nom
	image, _ := config["Image"].(string)
	name := c.Query("name")

	allowed, reason := filter.CheckContainerCreate(image, name, config)
	if !allowed {
		logger.Warnf("Container creation denied: %s", reason)
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Container creation denied by advanced filter",
			"reason":  reason,
		})
		c.Abort()
		return false
	}

	return true
}

// checkVolumeCreate vérifie la création de volume
func checkVolumeCreate(c *gin.Context, filter *filters.AdvancedFilter, logger *logrus.Logger) bool {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return false
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var config map[string]interface{}
	if err := json.Unmarshal(body, &config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return false
	}

	name, _ := config["Name"].(string)
	driver, _ := config["Driver"].(string)

	// Extraire le chemin du host si présent dans les options
	hostPath := ""
	if driverOpts, ok := config["DriverOpts"].(map[string]interface{}); ok {
		if device, ok := driverOpts["device"].(string); ok {
			hostPath = device
		}
		if path, ok := driverOpts["o"].(string); ok {
			if strings.Contains(path, "device=") {
				parts := strings.Split(path, "device=")
				if len(parts) > 1 {
					hostPath = strings.Split(parts[1], ",")[0]
				}
			}
		}
	}

	allowed, reason := filter.CheckVolumeMount(name, hostPath, driver)
	if !allowed {
		logger.Warnf("Volume creation denied: %s", reason)
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Volume creation denied by advanced filter",
			"reason":  reason,
		})
		c.Abort()
		return false
	}

	return true
}

// checkNetworkCreate vérifie la création de réseau
func checkNetworkCreate(c *gin.Context, filter *filters.AdvancedFilter, logger *logrus.Logger) bool {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return false
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var config map[string]interface{}
	if err := json.Unmarshal(body, &config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return false
	}

	name, _ := config["Name"].(string)
	driver, _ := config["Driver"].(string)

	allowed, reason := filter.CheckNetworkCreate(name, driver)
	if !allowed {
		logger.Warnf("Network creation denied: %s", reason)
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Network creation denied by advanced filter",
			"reason":  reason,
		})
		c.Abort()
		return false
	}

	return true
}

// checkImageCreate vérifie la création/pull d'image
func checkImageCreate(c *gin.Context, filter *filters.AdvancedFilter, logger *logrus.Logger) bool {
	fromImage := c.Query("fromImage")
	tag := c.Query("tag")

	imageName := fromImage
	if tag != "" {
		imageName += ":" + tag
	}

	allowed, reason := filter.CheckImageOperation(imageName)
	if !allowed {
		logger.Warnf("Image operation denied: %s", reason)
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Image operation denied by advanced filter",
			"reason":  reason,
		})
		c.Abort()
		return false
	}

	return true
}

// checkImageBuild vérifie la construction d'image
func checkImageBuild(c *gin.Context, filter *filters.AdvancedFilter, logger *logrus.Logger) bool {
	tag := c.Query("t")
	if tag == "" {
		// Pas de tag spécifié, on laisse passer
		return true
	}

	allowed, reason := filter.CheckImageOperation(tag)
	if !allowed {
		logger.Warnf("Image build denied: %s", reason)
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Image build denied by advanced filter",
			"reason":  reason,
		})
		c.Abort()
		return false
	}

	return true
}
