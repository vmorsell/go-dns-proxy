.PHONY: all
all: lint test build

.PHONY: lint
lint:
	staticcheck ./...

.PHONY: install-linter
install-linter:
	go get honnef.co/go/tools/cmd/staticcheck@latest

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build ./...

.PHONY: build-image
build-image:
	docker build --no-cache -t vmorsell/go-dns-proxy:latest .

.PHONY: run-image
run-image:
	docker run -p 53:53/tcp -p 53:53/udp vmorsell/go-dns-proxy
