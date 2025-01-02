.PHONY: all
all: build

.PHONY: build
build:
	go build .

.PHONY: install
install:
	go install .

.PHONY: lint
lint:
	go vet ./...

.PHONY: clean
clean:
	go clean
