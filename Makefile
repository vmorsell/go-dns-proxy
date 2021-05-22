.PHONY: build-container run-container

build-image:
	docker build --no-cache -t vmorsell/go-dns-proxy:latest .

run-image:
	docker run -p 53:53/tcp -p 53:53/udp vmorsell/go-dns-proxy
