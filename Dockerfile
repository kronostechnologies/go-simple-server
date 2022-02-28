FROM --platform=$BUILDPLATFORM golang:1.17 AS builder
WORKDIR /go/src/github.com/kronostechnologies/go-simple-server/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -ldflags="-w -s" -o simple-server .
RUN echo "nobody:x:65534:65534:nobody:/:" > /tmp/passwd

FROM scratch
COPY --from=builder /go/src/github.com/kronostechnologies/go-simple-server/simple-server /bin/
COPY --from=builder /tmp/passwd /etc/passwd
USER 65534:65534
ENTRYPOINT ["/bin/simple-server"]
