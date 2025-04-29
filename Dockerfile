# Build the application from source
FROM golang:1.24 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy all .go files recursively, excluding test files
COPY . .
RUN find . -type f -name "*_test.go" -delete

RUN CGO_ENABLED=0 GOOS=linux go build -o /dist/home-automation

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /dist/home-automation /usr/bin/home-automation

COPY example.config.yaml /etc/home-automation/config.yaml

ENV GIN_MODE=release

USER nonroot:nonroot

ENTRYPOINT ["/usr/bin/home-automation", "--config=/etc/home-automation/config.yaml"]