FROM golang:1.15 AS builder
WORKDIR /go/src/github.com/kronostechnologies/go-simple-server/
COPY * ./
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o simple-server .

FROM scratch
COPY --from=builder /go/src/github.com/kronostechnologies/go-simple-server/simple-server /bin/
ENTRYPOINT ["/bin/simple-server"]
