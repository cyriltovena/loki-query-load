PHONY: all build test clean build-image push-image
.DEFAULT_GOAL := all

IMAGE_PREFIX ?= ctovena
APP_NAME = loki-gen-load
IMAGE_TAG := 0.1
SRCS = $(wildcard *.go)

all: test build-image

build:
	go build -o $(APP_NAME) -v $(SRCS)

test:
	go test -v ./...

clean:
	rm -f $(APP_NAME)
	go clean ./...

build-image:
	docker build -t $(IMAGE_PREFIX)/$(APP_NAME) .
	docker tag $(IMAGE_PREFIX)/$(APP_NAME) $(IMAGE_PREFIX)/$(APP_NAME):$(IMAGE_TAG)

push-image:
	docker push $(IMAGE_PREFIX)/$(APP_NAME):$(IMAGE_TAG)
	docker push $(IMAGE_PREFIX)/$(APP_NAME):latest

deploy:
	kubectl apply -f deployment.yaml

delete:
	kubectl delete -f deployment.yaml
