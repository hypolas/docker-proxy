package proxy

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"docker-proxy/config"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

// Handler handles Docker API proxy requests
type Handler struct {
	client       *resty.Client
	dockerSocket string
	config       *config.Config
}

// NewHandler creates a new proxy handler
func NewHandler(cfg *config.Config) *Handler {
	client := resty.New()

	// Configure for Unix socket or TCP connection
	if strings.HasPrefix(cfg.DockerSocket, "unix://") {
		socketPath := strings.TrimPrefix(cfg.DockerSocket, "unix://")
		client.SetTransport(&http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		})
		client.SetBaseURL("http://unix")
	} else if strings.HasPrefix(cfg.DockerSocket, "tcp://") {
		client.SetBaseURL(cfg.DockerSocket)
	} else {
		// Default to unix socket
		client.SetTransport(&http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial("unix", cfg.DockerSocket)
			},
		})
		client.SetBaseURL("http://unix")
	}

	client.SetTimeout(60 * time.Second)
	client.SetRetryCount(0)

	return &Handler{
		client:       client,
		dockerSocket: cfg.DockerSocket,
		config:       cfg,
	}
}

// ProxyRequest proxies the request to Docker socket
func (h *Handler) ProxyRequest(c *gin.Context) {
	targetPath := buildTargetPath(c)
	req := h.prepareRequest(c)
	if req == nil {
		return // Error already handled
	}

	resp, err := h.executeRequest(req, targetPath, c)
	if err != nil {
		return // Error already handled
	}

	h.sendResponse(c, resp)
}

// buildTargetPath constructs the target path with query parameters
func buildTargetPath(c *gin.Context) string {
	targetPath := c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		targetPath += "?" + c.Request.URL.RawQuery
	}
	return targetPath
}

// prepareRequest creates and configures a resty request
func (h *Handler) prepareRequest(c *gin.Context) *resty.Request {
	req := h.client.R()

	// Copy headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.SetHeader(key, value)
		}
	}

	// Copy body if present
	if c.Request.Body != nil {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to read request body",
			})
			return nil
		}
		req.SetBody(body)
	}

	return req
}

// executeRequest executes the HTTP request using the appropriate method
func (h *Handler) executeRequest(req *resty.Request, targetPath string, c *gin.Context) (*resty.Response, error) {
	var resp *resty.Response
	var err error

	switch c.Request.Method {
	case http.MethodGet:
		resp, err = req.Get(targetPath)
	case http.MethodPost:
		resp, err = req.Post(targetPath)
	case http.MethodPut:
		resp, err = req.Put(targetPath)
	case http.MethodDelete:
		resp, err = req.Delete(targetPath)
	case http.MethodPatch:
		resp, err = req.Patch(targetPath)
	case http.MethodHead:
		resp, err = req.Head(targetPath)
	case http.MethodOptions:
		resp, err = req.Options(targetPath)
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": fmt.Sprintf("method %s not allowed", c.Request.Method),
		})
		return nil, fmt.Errorf("method not allowed")
	}

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": fmt.Sprintf("failed to proxy request: %v", err),
		})
		return nil, err
	}

	return resp, nil
}

// sendResponse sends the proxied response back to the client
func (h *Handler) sendResponse(c *gin.Context, resp *resty.Response) {
	// Copy response headers
	for key, values := range resp.Header() {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Send response
	c.Data(resp.StatusCode(), resp.Header().Get("Content-Type"), resp.Body())
}
