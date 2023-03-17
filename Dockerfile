FROM golang:1.18-alpine as builder

WORKDIR /app

# Copy the application source code to the container
COPY . .

# Compile the application
RUN go build -o carpool ./cmd/carpool/main.go

# Base image for the application container
FROM alpine:latest

# Get some basic things and remove unnecessary apk files
RUN apk --update upgrade && apk add \
  ca-certificates \
  curl \
  tzdata \
  bash \
  && update-ca-certificates \
  && rm -rf /var/cache/apk/*

# The port on which the service will listen
EXPOSE 8080

# Tag the image with a configurable label
ARG BUILD_TAG=unknown
LABEL BUILD_TAG=$BUILD_TAG

# Copy the compiled binary
COPY --from=builder /app /app
COPY --from=builder /app/carpool /app/carpool

# The command to execute
ENTRYPOINT ["/app/carpool"]
