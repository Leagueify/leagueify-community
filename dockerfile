# syntax=docker/dockerfile:1
# download go depedencies
FROM golang:1.22.0-alpine3.19 AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# build local development image
FROM golang:1.22.0-alpine3.19 AS dev
COPY --from=base /go/bin /go/bin
COPY --from=base /go/pkg /go/pkg
WORKDIR /app
RUN go install github.com/air-verse/air@latest
COPY . ./
EXPOSE 8888
# build the go binary
FROM golang:1.22.0-alpine3.19 AS builder
COPY --from=base /go/bin /go/bin
COPY --from=base /go/pkg /go/pkg
WORKDIR /app
COPY . ./
EXPOSE 8888
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/server cmd/server.go
# create production image
FROM gcr.io/distroless/base-debian11 AS release
COPY --from=builder /app/bin/server /server
EXPOSE 8888
ENTRYPOINT ["/server"]
