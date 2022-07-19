FROM golang:1.18-alpine AS buildenv
WORKDIR /src
ADD . /src
RUN go mod download
RUN go build -o sha256sum cmd/hasher/main.go

RUN chmod +x sha256sum

FROM alpine:latest
WORKDIR /app
VOLUME /app
COPY --from=buildenv /src/sha256sum .
COPY --from=buildenv /src/config.yaml ./
COPY --from=buildenv /src/.env .

ENTRYPOINT ["/app/sha256sum"]