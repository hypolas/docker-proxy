FROM  golang:1.24.7-trixie AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o docker-proxy ./cmd/docker-proxy

# Final stage
FROM alpine:latest as final

RUN apk --no-cache add ca-certificates sudo su-exec util-linux

# create user dkrproxy and add to sudoers and set password ad add to group sudo
RUN addgroup -S dkrproxy && \
    adduser -S -G dkrproxy dkrproxy && \
    echo "dkrproxy ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers


WORKDIR /app/

# Copy binary from builder
COPY --from=builder /app/docker-proxy .

# Copy entrypoint wrapper
COPY entrypoint.sh /entrypoint.sh

# Expose port
EXPOSE 2375

# Copy sudoers file
RUN cp -r /etc /etc_source

# Run the proxy by default via entrypoint script
# Stay as root to allow setpriv to work in entrypoint
ENTRYPOINT ["/entrypoint.sh"]
CMD ["docker-proxy"]
