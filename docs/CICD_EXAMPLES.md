# CI/CD Integration - Docker Socket Proxy

This guide details how to integrate the proxy into different CI/CD platforms to secure your Docker pipelines.

> ℹ️ Sur des runners auto-hébergés vous pouvez exposer le proxy via un socket Unix en définissant `LISTEN_SOCKET=/tmp/docker-proxy.sock` (et éventuellement `SOCKET_PERMS`). Les jobs doivent alors utiliser `DOCKER_HOST=unix:///tmp/docker-proxy.sock`.

## 🎯 Why Use This Proxy in CI/CD?

### Problems Solved
- ❌ **Unrestricted access to Docker socket** in CI/CD runners
- ❌ **Privilege escalation** via `docker run --privileged`
- ❌ **Docker socket mounting** in build containers
- ❌ **Use of unversioned :latest** tags
- ❌ **Pull from unapproved public registries**

### Benefits
- ✅ **Granular control**: Allow only build, push, pull operations
- ✅ **Enforced private registry**: Block Docker Hub, allow only your registry
- ✅ **Enforced versioning**: Forbid :latest, :dev, :test tags
- ✅ **Complete audit**: Logs of all Docker operations
- ✅ **Zero-trust**: Principle of least privilege applied

---

## 🔧 GitHub Actions

### Basic Configuration

```yaml
# .github/workflows/docker-build.yml
name: Secure Docker Build

on:
  push:
    branches: [main, develop]

services:
  docker-proxy:
    image: ghcr.io/your-org/docker-proxy:latest
    env:
      # Allowed endpoints
      CONTAINERS: "1"
      IMAGES: "1"
      BUILD: "1"
      POST: "1"

      # CI/CD security
      DKRPRX__CONTAINERS__ALLOWED_IMAGES: "^ghcr.io/your-org/.*"
      DKRPRX__IMAGES__DENIED_TAGS: "^latest$"
      DKRPRX__CONTAINERS__DENY_PRIVILEGED: "true"

      # Proxy protection
      PROXY_CONTAINER_NAME: docker-proxy
      LISTEN_SOCKET: unix:///tmp/docker-proxy.sock

    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - /tmp:/tmp

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Docker via Proxy
        run: |
          export DOCKER_HOST=unix:///tmp/docker-proxy.sock
          docker version

      - name: Build Docker image
        env:
          DOCKER_HOST: unix:///tmp/docker-proxy.sock
        run: |
          docker build -t ghcr.io/your-org/app:${{ github.sha }} .

      - name: Push to registry
        env:
          DOCKER_HOST: unix:///tmp/docker-proxy.sock
        run: |
          echo "${{ secrets.GITHUB_TOKEN }}" | docker login ghcr.io -u ${{ github.actor }} --password-stdin
          docker push ghcr.io/your-org/app:${{ github.sha }}
```

### Advanced Configuration (Multi-Environment)

```yaml
# .github/workflows/secure-deploy.yml
name: Multi-Environment Secure Deploy

on:
  push:
    branches: [main]

services:
  docker-proxy:
    image: ghcr.io/your-org/docker-proxy:latest
    env:
      CONTAINERS: "1"
      IMAGES: "1"
      BUILD: "1"
      NETWORKS: "1"
      POST: "1"
      DELETE: "1"

      # Advanced filters
      DKRPRX__CONTAINERS__ALLOWED_IMAGES: "^ghcr.io/your-org/.*:v[0-9]+\\.[0-9]+\\.[0-9]+$"
      DKRPRX__CONTAINERS__REQUIRE_LABELS: "ci=github-actions,team=platform"
      DKRPRX__IMAGES__DENIED_TAGS: "^(latest|dev|test)$"
      DKRPRX__CONTAINERS__DENY_HOST_NETWORK: "true"
      LISTEN_SOCKET: unix:///tmp/docker-proxy.sock

    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - /tmp:/tmp

jobs:
  build-production:
    runs-on: ubuntu-latest
    env:
      DOCKER_HOST: unix:///tmp/docker-proxy.sock
    steps:
      - uses: actions/checkout@v3

      - name: Build with required labels
        run: |
          docker build \
            --label ci=github-actions \
            --label team=platform \
            -t ghcr.io/your-org/app:v1.0.${{ github.run_number }} .
```

---

## 🦊 GitLab CI/CD

### Basic Configuration

```yaml
# .gitlab-ci.yml
variables:
  DOCKER_HOST: unix:///tmp/docker-proxy.sock
  DOCKER_DRIVER: overlay2

services:
  - name: registry.gitlab.com/your-org/docker-proxy:latest
    alias: docker-proxy
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - /tmp:/tmp
    variables:
      CONTAINERS: "1"
      IMAGES: "1"
      BUILD: "1"
      POST: "1"
      DELETE: "1"
      LISTEN_SOCKET: unix:///tmp/docker-proxy.sock

      # Force GitLab registry only
      DKRPRX__CONTAINERS__ALLOWED_IMAGES: "^registry.gitlab.com/your-org/.*"
      DKRPRX__IMAGES__ALLOWED_REPOS: "^registry.gitlab.com/your-org/.*"

      # Forbid unversioned tags
      DKRPRX__IMAGES__DENIED_TAGS: "^(latest|dev|test|master|main)$"

stages:
  - build
  - test
  - deploy

build:
  stage: build
  script:
    - docker build -t registry.gitlab.com/your-org/app:$CI_COMMIT_SHA .
    - docker push registry.gitlab.com/your-org/app:$CI_COMMIT_SHA
  only:
    - branches

build-release:
  stage: build
  script:
    - docker build -t registry.gitlab.com/your-org/app:$CI_COMMIT_TAG .
    - docker push registry.gitlab.com/your-org/app:$CI_COMMIT_TAG
  only:
    - tags
```

### Configuration with Cache and Multi-Stage

```yaml
# .gitlab-ci.yml with Docker cache
.docker-proxy:
  services:
    - name: registry.gitlab.com/your-org/docker-proxy:latest
      alias: docker-proxy
      volumes:
        - /var/run/docker.sock:/var/run/docker.sock:ro
        - /tmp:/tmp
      variables:
        CONTAINERS: "1"
        IMAGES: "1"
        BUILD: "1"
        POST: "1"
        LISTEN_SOCKET: unix:///tmp/docker-proxy.sock
        DKRPRX__CONTAINERS__ALLOWED_IMAGES: "^registry.gitlab.com/your-org/.*"
        DKRPRX__IMAGES__DENIED_TAGS: "^latest$"

build-optimized:
  extends: .docker-proxy
  stage: build
  variables:
    DOCKER_HOST: unix:///tmp/docker-proxy.sock
  script:
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
    - docker build --cache-from registry.gitlab.com/your-org/app:cache -t registry.gitlab.com/your-org/app:$CI_COMMIT_SHA .
    - docker push registry.gitlab.com/your-org/app:$CI_COMMIT_SHA
```

---

## 🔵 Azure DevOps

### Pipeline Configuration

```yaml
# azure-pipelines.yml
trigger:
  - main

pool:
  vmImage: 'ubuntu-latest'

resources:
  containers:
    - container: docker-proxy
      image: yourregistry.azurecr.io/docker-proxy:latest
      env:
        CONTAINERS: 1
        IMAGES: 1
        BUILD: 1
        POST: 1
        LISTEN_SOCKET: unix:///tmp/docker-proxy.sock
        DKRPRX__CONTAINERS__ALLOWED_IMAGES: '^yourregistry\.azurecr\.io/.*'
        DKRPRX__IMAGES__DENIED_TAGS: '^latest$'
      volumes:
        - /var/run/docker.sock:/var/run/docker.sock:ro
        - /tmp:/tmp

services:
  docker-proxy: docker-proxy

steps:
  - script: |
      export DOCKER_HOST=unix:///tmp/docker-proxy.sock
      docker build -t yourregistry.azurecr.io/app:$(Build.BuildId) .
    displayName: 'Build Docker Image'

  - script: |
      export DOCKER_HOST=unix:///tmp/docker-proxy.sock
      docker push yourregistry.azurecr.io/app:$(Build.BuildId)
    displayName: 'Push to ACR'
```

---

## 🟢 Jenkins

### Jenkinsfile with Docker Proxy

```groovy
// Jenkinsfile
pipeline {
    agent any

    environment {
        DOCKER_HOST = 'unix:///tmp/docker-proxy.sock'
        REGISTRY = 'registry.company.com'
    }

    stages {
        stage('Start Docker Proxy') {
            steps {
                sh '''
                    docker run -d \
                      --name docker-proxy-${BUILD_NUMBER} \
                      -v /var/run/docker.sock:/var/run/docker.sock:ro \
                      -v /tmp:/tmp \
                      -e LISTEN_SOCKET=unix:///tmp/docker-proxy.sock \
                      -e CONTAINERS=1 \
                      -e IMAGES=1 \
                      -e BUILD=1 \
                      -e POST=1 \
                      -e DKRPRX__CONTAINERS__ALLOWED_IMAGES="^${REGISTRY}/.*" \
                      -e DKRPRX__IMAGES__DENIED_TAGS="^latest$" \
                      -e PROXY_CONTAINER_NAME=docker-proxy-${BUILD_NUMBER} \
                      registry.company.com/docker-proxy:latest
                '''
            }
        }

        stage('Build') {
            steps {
                sh '''
                    docker build -t ${REGISTRY}/app:${BUILD_NUMBER} .
                '''
            }
        }

        stage('Push') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'docker-registry', usernameVariable: 'USER', passwordVariable: 'PASS')]) {
                    sh '''
                        echo $PASS | docker login ${REGISTRY} -u $USER --password-stdin
                        docker push ${REGISTRY}/app:${BUILD_NUMBER}
                    '''
                }
            }
        }
    }

    post {
        always {
            sh 'docker rm -f docker-proxy-${BUILD_NUMBER} || true'
        }
    }
}
```

---

## 🔴 CircleCI

### Configuration .circleci/config.yml

```yaml
version: 2.1

executors:
  docker-proxy:
    docker:
      - image: cimg/base:stable
      - image: hypolas/proxy-docker:latest
        name: docker-proxy
        volumes:
          - /var/run/docker.sock:/var/run/docker.sock:ro
          - /tmp:/tmp
        environment:
          CONTAINERS: "1"
          IMAGES: "1"
          BUILD: "1"
          POST: "1"
          LISTEN_SOCKET: unix:///tmp/docker-proxy.sock
          DKRPRX__CONTAINERS__ALLOWED_IMAGES: "^your-registry\\.io/.*"
          DKRPRX__IMAGES__DENIED_TAGS: "^latest$"

jobs:
  build-and-push:
    executor: docker-proxy
    steps:
      - checkout
      - setup_remote_docker:
          docker_layer_caching: true

      - run:
          name: Build Image
          environment:
            DOCKER_HOST: unix:///tmp/docker-proxy.sock
          command: |
            docker build -t your-registry.io/app:${CIRCLE_SHA1} .

      - run:
          name: Push Image
          environment:
            DOCKER_HOST: unix:///tmp/docker-proxy.sock
          command: |
            echo $DOCKER_PASSWORD | docker login -u $DOCKER_USERNAME --password-stdin your-registry.io
            docker push your-registry.io/app:${CIRCLE_SHA1}

workflows:
  build-deploy:
    jobs:
      - build-and-push
```

---

## 📋 Secure CI/CD Checklist

- [ ] **Docker socket read-only**: `:ro` on all mounts
- [ ] **Private registry enforced**: `DKRPRX__CONTAINERS__ALLOWED_IMAGES`
- [ ] **Forbid :latest**: `DKRPRX__IMAGES__DENIED_TAGS="^latest$"`
- [ ] **Read-only by default**: `POST=0` except build/push
- [ ] **Forbid privileged**: `DKRPRX__CONTAINERS__DENY_PRIVILEGED=true`
- [ ] **Forbid host network**: `DKRPRX__CONTAINERS__DENY_HOST_NETWORK=true`
- [ ] **Required labels**: `DKRPRX__CONTAINERS__REQUIRE_LABELS`
- [ ] **Proxy protection**: `PROXY_CONTAINER_NAME` configured
- [ ] **Centralized logs**: Collect and analyze proxy logs
- [ ] **Regular audit**: Check blocked attempts

---

## 🔍 CI/CD Troubleshooting

### Image Rejected
```
Error: Container creation denied by advanced filter
Reason: image not in allowed list: docker.io/nginx:latest
```
**Solution**: Check `DKRPRX__CONTAINERS__ALLOWED_IMAGES`

### Tag Rejected
```
Error: Image operation denied by advanced filter
Reason: image tag is denied: latest
```
**Solution**: Use a versioned tag: `v1.0.0`

### Privileged Container Rejected
```
Error: Container creation denied by advanced filter
Reason: privileged containers are denied
```
**Solution**: Remove `--privileged` or set `DKRPRX__CONTAINERS__DENY_PRIVILEGED=false`

---

## 📚 Resources

- [Complete Documentation](../README.md)
- [Advanced Filters](../docs/ADVANCED_FILTERS.md)
- [Environment Variable Configuration](../docs/ENV_FILTERS.md)
- [Security Guide](../docs/SECURITY.md)
