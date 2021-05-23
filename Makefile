.PHONY: all
all: lint test build

.PHONY: lint
lint:
	staticcheck ./...

.PHONY: get-linter
get-linter:
	go get honnef.co/go/tools/cmd/staticcheck@latest

.PHONY: test
test:
	go test ./...

build:
	go build ./...

build-image:
	docker build --no-cache -t vmorsell/go-dns-proxy:latest .

run-image:
	docker run -p 53:53/tcp -p 53:53/udp vmorsell/go-dns-proxy
