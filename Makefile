BASE_PATH = $(PWD)
GOSRC_PATH = $(PWD)/src
OUTPUT_PATH = $(PWD)/build
OUTPUT_NAME = version_monitor
GOBIN = go
DOCKERBIN = docker
GOBUILD = $(GOBIN) build
BUILD_SHA = $(shell git rev-parse HEAD)
SOURCES = $(subst $(GOSRC_PATH)/,,$(wildcard $(GOSRC_PATH)/*.go))
LDFLAGS = -ldflags "-X main.BuildSHA=$(BUILD_SHA)"
DOCKER_TAG = vmonitor:latest
DOCKER_IMAGE = $(OUTPUT_NAME).tar

all: build

clean:
	rm -rf $(OUTPUT_PATH)

$(OUTPUT_PATH)/$(OUTPUT_NAME):
	$(DOCKERBIN) run --rm -w /usr/src \
    -e GO15VENDOREXPERIMENT="1" \
    -v $(OUTPUT_PATH):/usr/src/build \
    -v $(GOSRC_PATH):/usr/src \
    -v $(GOSRC_PATH)/vendor:/go/src \
    golang:1.5.1-alpine \
    $(GOBUILD) $(LDFLAGS) -v -o /usr/src/build/$(OUTPUT_NAME) $(SOURCES)

$(OUTPUT_PATH)/config.yml:
	cp config.yaml $(OUTPUT_PATH)/config.yml

build: $(OUTPUT_PATH)/$(OUTPUT_NAME) $(OUTPUT_PATH)/config.yml
	$(DOCKERBIN) build -t $(DOCKER_TAG) $(BASE_PATH)

docker-image:
	$(DOCKERBIN) save -o $(DOCKER_IMAGE) $(DOCKER_TAG)

.PHONY: build clean docker-image
