package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"docker-proxy/config"
	"docker-proxy/internal/middleware"
	"docker-proxy/internal/proxy"
	"docker-proxy/pkg/rules"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup logger
	logger := setupLogger(cfg.LogLevel)
	logger.Info("Starting Docker Socket Proxy")

	// Auto-detect Docker API version if not set
	if cfg.APIVersion == "" {
		cfg.APIVersion = config.DetectDockerAPIVersion(cfg.DockerSocket, logger)
	} else {
		logger.Infof("Using configured Docker API version: %s", cfg.APIVersion)
	}

	// Create rule matcher
	matcher := rules.NewMatcher(cfg.AccessRules)

	// Create proxy handler
	proxyHandler := proxy.NewHandler(cfg)

	// Setup Gin router
	if cfg.LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.ACLMiddleware(matcher))
	router.Use(middleware.AdvancedFilterMiddleware(cfg.AdvancedFilters, logger))

	// Register catch-all route for proxying
	router.Any("/*path", proxyHandler.ProxyRequest)

	// Log configuration
	logConfiguration(logger, cfg)

	// Create HTTP server
	srv := &http.Server{
		Handler: router,
	}

	// Start server (Unix socket or TCP)
	go func() {
		if cfg.ListenSocket != "" {
			// Unix socket mode
			// Strip "unix://" prefix if present
			socketPath := cfg.ListenSocket
			if len(socketPath) > 7 && socketPath[:7] == "unix://" {
				socketPath = socketPath[7:]
			}

			logger.Infof("Listening on Unix socket: %s", socketPath)
			logger.Infof("Proxying to %s", cfg.DockerSocket)

			// Remove existing socket file if it exists
			os.Remove(socketPath)

			listener, err := net.Listen("unix", socketPath)
			if err != nil {
				logger.Fatalf("Failed to create Unix socket: %v", err)
			}

			// Set socket permissions (0666 for wider access, or use SOCKET_PERMS env var)
			socketPerms := os.Getenv("SOCKET_PERMS")
			if socketPerms == "" {
				socketPerms = "0666" // Default: accessible by all users
			}

			// Parse octal permission string
			perms := os.FileMode(0666) // default
			if permValue, err := strconv.ParseUint(socketPerms, 8, 32); err != nil {
				logger.Warnf("Invalid SOCKET_PERMS '%s', using 0666", socketPerms)
			} else {
				perms = os.FileMode(permValue)
			}

			if err := os.Chmod(socketPath, perms); err != nil {
				logger.Warnf("Failed to set socket permissions: %v", err)
			} else {
				logger.Infof("Socket permissions set to %s", socketPerms)
			}

			if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
				logger.Fatalf("Failed to start server: %v", err)
			}
		} else {
			// TCP mode
			logger.Infof("Listening on %s", cfg.ListenAddr)
			logger.Infof("Proxying to %s", cfg.DockerSocket)

			srv.Addr = cfg.ListenAddr
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Fatalf("Failed to start server: %v", err)
			}
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	// Cleanup Unix socket if used
	if cfg.ListenSocket != "" {
		socketPath := cfg.ListenSocket
		if len(socketPath) > 7 && socketPath[:7] == "unix://" {
			socketPath = socketPath[7:]
		}
		os.Remove(socketPath)
	}

	logger.Info("Server stopped")
}

// setupLogger configures the logger based on log level
func setupLogger(level string) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger
}

// logConfiguration logs the current access rules configuration
func logConfiguration(logger *logrus.Logger, cfg *config.Config) {
	logger.Info("Access Rules Configuration:")

	rules := cfg.AccessRules

	// Log granted endpoints
	granted := []string{}
	if rules.Events {
		granted = append(granted, "EVENTS")
	}
	if rules.Ping {
		granted = append(granted, "PING")
	}
	if rules.Version {
		granted = append(granted, "VERSION")
	}
	if rules.Auth {
		granted = append(granted, "AUTH")
	}
	if rules.Build {
		granted = append(granted, "BUILD")
	}
	if rules.Commit {
		granted = append(granted, "COMMIT")
	}
	if rules.Configs {
		granted = append(granted, "CONFIGS")
	}
	if rules.Containers {
		granted = append(granted, "CONTAINERS")
	}
	if rules.Distribution {
		granted = append(granted, "DISTRIBUTION")
	}
	if rules.Exec {
		granted = append(granted, "EXEC")
	}
	if rules.Images {
		granted = append(granted, "IMAGES")
	}
	if rules.Info {
		granted = append(granted, "INFO")
	}
	if rules.Networks {
		granted = append(granted, "NETWORKS")
	}
	if rules.Nodes {
		granted = append(granted, "NODES")
	}
	if rules.Plugins {
		granted = append(granted, "PLUGINS")
	}
	if rules.Secrets {
		granted = append(granted, "SECRETS")
	}
	if rules.Services {
		granted = append(granted, "SERVICES")
	}
	if rules.Session {
		granted = append(granted, "SESSION")
	}
	if rules.Swarm {
		granted = append(granted, "SWARM")
	}
	if rules.System {
		granted = append(granted, "SYSTEM")
	}
	if rules.Tasks {
		granted = append(granted, "TASKS")
	}
	if rules.Volumes {
		granted = append(granted, "VOLUMES")
	}

	logger.Infof("  Granted endpoints: %v", granted)

	// Log HTTP methods
	methods := []string{"GET", "HEAD"}
	if rules.Post {
		methods = append(methods, "POST")
	}
	if rules.Delete {
		methods = append(methods, "DELETE")
	}
	if rules.Put {
		methods = append(methods, "PUT", "PATCH")
	}

	logger.Infof("  Allowed methods: %v", methods)

	if !rules.Post && !rules.Delete && !rules.Put {
		logger.Warn("  ⚠️  Read-only mode enabled (POST, DELETE, PUT disabled)")
	}
}
