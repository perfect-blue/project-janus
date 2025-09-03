BINARY_NAME=janus
IMAGE_NAME=janus:latest
.DEFAULT_GOAL := build

build:
	go build -o $(BINARY_NAME) ./src

run: build
	./$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)

docker-build:
	docker build -t $(IMAGE_NAME) .

docker-run:
	docker run --rm -p 8080:8080 $(IMAGE_NAME)

up: docker-build docker-run