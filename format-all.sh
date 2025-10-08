#!/bin/bash
# Format all Go files with gofmt

echo "Formatting all Go files..."
gofmt -w cmd/dockershield/main.go
gofmt -w config/config.go
gofmt -w config/env_filters.go
gofmt -w config/defaults.go
gofmt -w config/merge.go
gofmt -w config/version.go
gofmt -w config/config_test.go
gofmt -w internal/middleware/acl.go
gofmt -w internal/middleware/advanced_filter.go
gofmt -w internal/middleware/logging.go
gofmt -w internal/middleware/acl_test.go
gofmt -w internal/proxy/handler.go
gofmt -w pkg/filters/advanced.go
gofmt -w pkg/filters/json.go
gofmt -w pkg/filters/advanced_test.go
gofmt -w pkg/rules/matcher.go
gofmt -w pkg/rules/matcher_test.go
echo "Done!"
