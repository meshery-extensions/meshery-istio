# Copyright Meshery Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

include build/Makefile.core.mk
include build/Makefile.show-help.mk

#-----------------------------------------------------------------------------
# Environment Setup
#-----------------------------------------------------------------------------
BUILDER=buildx-multi-arch
ADAPTER=istio

#-----------------------------------------------------------------------------
# Docker-based Builds
#-----------------------------------------------------------------------------
.PHONY: docker docker-run lint error test run run-force-dynamic-reg


## Lint check Golang
lint:
	golangci-lint run

## Build Adapter container image with "edge-latest" tag
docker:
	DOCKER_BUILDKIT=1 docker build -t meshery/meshery-$(ADAPTER):$(RELEASE_CHANNEL)-latest .

## Run Adapter container with "edge-latest" tag
docker-run:
	(docker rm -f meshery-$(ADAPTER)) || true
	docker run --name meshery-$(ADAPTER) -d \
	-p 10000:10000 \
	-e DEBUG=true \
	meshery/meshery-$(ADAPTER):$(RELEASE_CHANNEL)-latest

## Build and run Adapter locally
run: dep-check
	go mod tidy; \
	DEBUG=true GOPROXY=direct GOSUMDB=off go run main.go

## Build and run Adapter locally; force component registration
run-force-dynamic-reg: dep-check
	FORCE_DYNAMIC_REG=true DEBUG=true GOPROXY=direct GOSUMDB=off go run main.go

## Run Meshery Error utility
error: dep-check
	go run github.com/meshery/meshkit/cmd/errorutil -d . analyze -i ./helpers -o ./helpers

## Run Golang tests
test:
	export CURRENTCONTEXT="$(kubectl config current-context)" 
	echo "current-context:" ${CURRENTCONTEXT} 
	export KUBECONFIG="${HOME}/.kube/config"
	echo "environment-kubeconfig:" ${KUBECONFIG}
	GOPROXY=direct GOSUMDB=off GO111MODULE=on go test -v ./...

#-----------------------------------------------------------------------------
# Dependencies
#-----------------------------------------------------------------------------
.PHONY: dep-check
#.SILENT: dep-check

INSTALLED_GO_VERSION=$(shell go version)

dep-check:

ifeq (,$(findstring $(GOVERSION), $(INSTALLED_GO_VERSION)))
# Only send a warning.
	@echo "Dependency missing: go$(GOVERSION). Ensure 'go$(GOVERSION).x' is installed and available in your 'PATH'"
	@echo "GOVERSION: " $(GOVERSION)
	@echo "INSTALLED_GO_VERSION: " $(INSTALLED_GO_VERSION)
# Force error and stop.
#	$(error Found $(INSTALLED_GO_VERSION). \
#	 Required golang version is: 'go$(GOVERSION).x'. \
#	 Ensure go '$(GOVERSION).x' is installed and available in your 'PATH'.)
endif
