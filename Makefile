VERSION = $(shell cat version.txt)
CACHE_PATH = cache
BLOCKLIST_USER = test
BLOCKLIST_PASS = test
DB_PATH = data/blocklist.db
BINARY_NAME = blocklist-cli
PORT = 3000
BIN_PATH = "bin/$(BINARY_NAME)"
DOMAIN=gcr.io
NAMESPACE=securework_homework



generate-schema:
	go run github.com/99designs/gqlgen generate


build:
	go build -o $(BIN_PATH) -ldflags "-X main.Version=$(VERSION)" *.go

build-all: generate-schema build

build-docker-image:
	docker build -t $(DOMAIN)/$(NAMESPACE)/${BINARY_NAME}:$(VERSION) \
		--build-arg PORT=$(PORT) \
		--build-arg CACHE_PATH=$(CACHE_PATH) \
		--build-arg BLOCKLIST_USER=$(BLOCKLIST_USER) \
		--build-arg BLOCKLIST_PASS=$(BLOCKLIST_PASS) \
		--build-arg DB_PATH=$(DB_PATH) .

publish-docker-image:
	docker push $(DOMAIN)/$(NAMESPACE)/${BINARY_NAME}:$(VERSION)

docker-build-and-push: build-docker-image publish-docker-image