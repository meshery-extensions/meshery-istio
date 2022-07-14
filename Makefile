GOPATH = $(shell go env GOPATH)
BUILDER=buildx-multi-arch

GIT_VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1`)
GIT_STRIPPED_VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1` | cut -c 2-)
v ?= 1.17.8 # Default go version to be used

check:
	golangci-lint run

protoc-setup:
	cd meshes
	wget https://raw.githubusercontent.com/layer5io/meshery/master/meshes/meshops.proto

proto:	
	protoc -I meshes/ meshes/meshops.proto --go_out=plugins=grpc:./meshes/

docker:
	DOCKER_BUILDKIT=1 docker build -t layer5/meshery-istio .

docker-run:
	(docker rm -f meshery-istio) || true
	docker run --name meshery-istio -d \
	-p 10000:10000 \
	-e DEBUG=true \
	layer5/meshery-istio:edge-latest

run:
	go$(v) mod tidy -compat=1.17; \
	DEBUG=true GOPROXY=direct GOSUMDB=off go run main.go

run-force-dynamic-reg:
	FORCE_DYNAMIC_REG=true DEBUG=true GOPROXY=direct GOSUMDB=off go run main.go
error:
	go run github.com/layer5io/meshkit/cmd/errorutil -d . analyze -i ./helpers -o ./helpers

test:
	export CURRENTCONTEXT="$(kubectl config current-context)" 
	echo "current-context:" ${CURRENTCONTEXT} 
	export KUBECONFIG="${HOME}/.kube/config"
	echo "environment-kubeconfig:" ${KUBECONFIG}
	GOPROXY=direct GOSUMDB=off GO111MODULE=on go test -v ./...
