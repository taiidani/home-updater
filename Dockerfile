FROM golang:1.15.8-alpine

# Build the app, dependencies first
RUN apk add --no-cache git
COPY go.mod go.sum /app/
WORKDIR /app
RUN go mod download

COPY . /app
ENV CGO_ENABLED=0
RUN go build -o main
RUN go test ./...

# ---
FROM alpine:3.13.2 AS dist

# Dependencies
RUN apk add --no-cache ca-certificates

# Add pre-built application
COPY --from=0 /app/main /app

ENTRYPOINT [ "/app" ]