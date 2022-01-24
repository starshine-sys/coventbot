FROM golang:latest AS builder

WORKDIR /build
COPY . ./
RUN go mod download
ENV CGO_ENABLED 0
RUN go build -v -o tribble -ldflags="-X github.com/starshine-sys/tribble/commands/static/info.GitVer=`git rev-parse --short HEAD`" ./cmd/tribble/

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /build/tribble tribble

CMD ["/app/tribble"]
