BASE_PATH = $(PWD)
GOSRC_PATH = $(PWD)/src
OUTPUT_PATH = $(PWD)/build
OUTPUT_NAME = version_monitor
VERSION = 0.0.1
GO_BIN = go
BUILD_TIME = $(shell TZ=utc date +%FT%T%z)
SOURCES = $(wildcard $(GOSRC_PATH)/*.go)
APP_NAME = Version_monitor
LDFLAGS = -ldflags "-X main.AppName=$(APP_NAME) -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

all: build

clean:
	rm -rf $(OUTPUT_PATH)

build-go:
	$(GO_BIN) build $(LDFLAGS) -o $(OUTPUT_PATH)/$(OUTPUT_NAME) $(SOURCES)

build: build-go

.PHONY: build build-go clean
