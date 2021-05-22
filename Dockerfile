FROM golang:1.16 AS builder
WORKDIR /go/src/github.com/vmorsell/go-dns-proxy/
RUN go get -d -v github.com/vmorsell/go-dns-proxy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/vmorsell/go-dns-proxy/main .
CMD ["./main"]
