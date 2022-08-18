TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=registry.terraform.io
NAMESPACE=malachantrio
NAME=eventstore
BINARY=terraform-provider-${NAME}
VERSION=0.0.1-dev
OS_ARCH=darwin_amd64

export GOPRIVATE=github.com/madedotcom/eventstore-client-go

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
