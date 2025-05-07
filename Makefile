# Makefile

BINARY_NAME := $(shell basename `pwd`)

BUILD_IMAGE := ghcr.io/riete/golang:1.23.6-busybox
ROOT := /$(shell basename `pwd`)


.PHONY: build-amd64
build-amd64:
	docker run --rm -w $(ROOT) --platform linux/amd64 -v .:$(ROOT) $(BUILD_IMAGE) go build -ldflags="-s -w" -o $(BINARY_NAME)-amd64 cmd/main.go
	upx $(BINARY_NAME)-amd64


.PHONY: build-arm64
build-arm64:
	docker run --rm -w $(ROOT) --platform linux/arm64 -v .:$(ROOT) $(BUILD_IMAGE) go build -ldflags="-s -w" -o $(BINARY_NAME)-arm64 cmd/main.go
	upx $(BINARY_NAME)-arm64


.PHONY: build
build: build-amd64 build-arm64
