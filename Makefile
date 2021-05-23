LINT_VERSION := 2020.2.3
LINT_CHECKSUM := "03b100561e3bc14db0b3b4004b102a00cb0197938d23cc40193f269f7b246d2d  staticcheck_linux_amd64.tar.gz"

.PHONY: all
all: lint test build

.PHONY: lint
lint:
	staticcheck ./...

.PHONY: install-linter
install-linter:
	go get honnef.co/go/tools/cmd/staticcheck@${LINT_VERSION}

.PHONY: install-linter-ci
install-linter-ci:
	wget https://github.com/dominikh/go-tools/releases/download/$(LINT_VERSION)/staticcheck_linux_amd64.tar.gz
	echo ${LINT_CHECKSUM} > checksums.txt
	sha256sum -c checksums.txt || exit 1
	tar xf staticcheck_linux_amd64.tar.gz staticcheck/staticcheck
	sudo cp staticcheck/staticcheck /usr/local/bin
	rm -rf staticcheck staticcheck_linux_amd64.tar.gz

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
