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

	"dockershield/config"
	"dockershield/internal/middleware"
	"dockershield/internal/proxy"
	"dockershield/pkg/rules"

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
	// Advanced filters run FIRST to allow DKRPRX__ variables to override ACL
	router.Use(middleware.AdvancedFilterMiddleware(cfg.AdvancedFilters, logger))
	router.Use(middleware.ACLMiddleware(matcher))

	// Register catch-all route for proxying
	router.Any("/*path", proxyHandler.ProxyRequest)

	// Log configuration
	logConfiguration(logger, cfg)

	// Create HTTP server
	srv := &http.Server{
		Handler:           router,
		ReadHeaderTimeout: 30 * time.Second, // Prevent Slowloris attacks
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Start server (Unix socket or TCP)
	go startServer(srv, cfg, logger)

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
		os.Remove(stripUnixPrefix(cfg.ListenSocket))
	}

	logger.Info("Server stopped")
}

// startServer starts the HTTP server on either Unix socket or TCP
func startServer(srv *http.Server, cfg *config.Config, logger *logrus.Logger) {
	if cfg.ListenSocket != "" {
		startUnixSocketServer(srv, cfg, logger)
	} else {
		startTCPServer(srv, cfg, logger)
	}
}

// startUnixSocketServer starts the server on a Unix socket
func startUnixSocketServer(srv *http.Server, cfg *config.Config, logger *logrus.Logger) {
	socketPath := stripUnixPrefix(cfg.ListenSocket)

	logger.Infof("Listening on Unix socket: %s", socketPath)
	logger.Infof("Proxying to %s", cfg.DockerSocket)

	// Remove existing socket file if it exists
	os.Remove(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		logger.Fatalf("Failed to create Unix socket: %v", err)
	}

	// Set socket permissions
	setSocketPermissions(socketPath, logger)

	if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Failed to start server: %v", err)
	}
}

// startTCPServer starts the server on a TCP address
func startTCPServer(srv *http.Server, cfg *config.Config, logger *logrus.Logger) {
	logger.Infof("Listening on %s", cfg.ListenAddr)
	logger.Infof("Proxying to %s", cfg.DockerSocket)

	srv.Addr = cfg.ListenAddr
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Failed to start server: %v", err)
	}
}

// setSocketPermissions sets permissions on a Unix socket
func setSocketPermissions(socketPath string, logger *logrus.Logger) {
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
}

// stripUnixPrefix removes "unix://" prefix from socket path
func stripUnixPrefix(socketPath string) string {
	if len(socketPath) > 7 && socketPath[:7] == "unix://" {
		return socketPath[7:]
	}
	return socketPath
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
	logger.Infof("  Granted endpoints: %v", getGrantedEndpoints(*cfg.AccessRules))
	logger.Infof("  Allowed methods: %v", getAllowedMethods(*cfg.AccessRules))

	if !cfg.AccessRules.Post && !cfg.AccessRules.Delete && !cfg.AccessRules.Put {
		logger.Warn("  ⚠️  Read-only mode enabled (POST, DELETE, PUT disabled)")
	}
}

// getGrantedEndpoints returns a list of enabled endpoints
func getGrantedEndpoints(rules config.AccessRules) []string {
	endpoints := []struct {
		enabled bool
		name    string
	}{
		{rules.Events, "EVENTS"},
		{rules.Ping, "PING"},
		{rules.Version, "VERSION"},
		{rules.Auth, "AUTH"},
		{rules.Build, "BUILD"},
		{rules.Commit, "COMMIT"},
		{rules.Configs, "CONFIGS"},
		{rules.Containers, "CONTAINERS"},
		{rules.Distribution, "DISTRIBUTION"},
		{rules.Exec, "EXEC"},
		{rules.Images, "IMAGES"},
		{rules.Info, "INFO"},
		{rules.Networks, "NETWORKS"},
		{rules.Nodes, "NODES"},
		{rules.Plugins, "PLUGINS"},
		{rules.Secrets, "SECRETS"},
		{rules.Services, "SERVICES"},
		{rules.Session, "SESSION"},
		{rules.Swarm, "SWARM"},
		{rules.System, "SYSTEM"},
		{rules.Tasks, "TASKS"},
		{rules.Volumes, "VOLUMES"},
	}

	granted := []string{}
	for _, ep := range endpoints {
		if ep.enabled {
			granted = append(granted, ep.name)
		}
	}
	return granted
}

// getAllowedMethods returns a list of allowed HTTP methods
func getAllowedMethods(rules config.AccessRules) []string {
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
	return methods
}
