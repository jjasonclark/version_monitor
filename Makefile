BASE_PATH = $(PWD)
GOSRC_PATH = $(PWD)/src
OUTPUT_PATH = $(PWD)/build
OUTPUT_NAME = version_monitor
GOBIN = go
GOBUILD = $(GOBIN) build
BUILD_SHA = $(shell git rev-parse HEAD)
SOURCES = $(wildcard $(GOSRC_PATH)/*.go)
LDFLAGS = -ldflags "-X main.BuildSHA=$(BUILD_SHA)"
GO15VENDOREXPERIMENT = 1

all: build

clean:
	rm -rf $(OUTPUT_PATH)

build-go:
	$(GOBUILD) $(LDFLAGS) -o $(OUTPUT_PATH)/$(OUTPUT_NAME) $(SOURCES)

build: build-go

.PHONY: build build-go clean
