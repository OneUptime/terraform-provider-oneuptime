# Terraform Provider Makefile

HOSTNAME=registry.terraform.io
NAMESPACE=oneuptime
NAME=oneuptime
BINARY=terraform-provider-${NAME}
VERSION=12.0.11
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)

default: install

build:
	go build -o ${BINARY}

release:
	mkdir -p ./builds
	GOOS=darwin GOARCH=amd64 go build -o ./builds/${BINARY}_darwin_amd64
	GOOS=darwin GOARCH=arm64 go build -o ./builds/${BINARY}_darwin_arm64

	GOOS=freebsd GOARCH=386 go build -o ./builds/${BINARY}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./builds/${BINARY}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./builds/${BINARY}_freebsd_arm
	GOOS=freebsd GOARCH=arm64 go build -o ./builds/${BINARY}_freebsd_arm64

	GOOS=linux GOARCH=386 go build -o ./builds/${BINARY}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./builds/${BINARY}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./builds/${BINARY}_linux_arm
	GOOS=linux GOARCH=arm64 go build -o ./builds/${BINARY}_linux_arm64

	GOOS=openbsd GOARCH=386 go build -o ./builds/${BINARY}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./builds/${BINARY}_openbsd_amd64
	GOOS=openbsd GOARCH=arm go build -o ./builds/${BINARY}_openbsd_arm
	GOOS=openbsd GOARCH=arm64 go build -o ./builds/${BINARY}_openbsd_arm64

	GOOS=solaris GOARCH=amd64 go build -o ./builds/${BINARY}_solaris_amd64

	GOOS=windows GOARCH=386 go build -o ./builds/${BINARY}_windows_386.exe
	GOOS=windows GOARCH=amd64 go build -o ./builds/${BINARY}_windows_amd64.exe
	GOOS=windows GOARCH=arm go build -o ./builds/${BINARY}_windows_arm.exe
	GOOS=windows GOARCH=arm64 go build -o ./builds/${BINARY}_windows_arm64.exe

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	go test ./... -v $(TESTARGS) -timeout 120m

testacc:
	TF_ACC=1 go test $(go list ./... | grep -v examples) -v $(TESTARGS) -timeout 120m

clean:
	rm -f ${BINARY}
	rm -rf ./builds/*

.PHONY: build release install test testacc clean