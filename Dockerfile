# Stage 1: Build the Go application
FROM golang:1.25.5-alpine3.22 AS builder

# Set the working directory inside the container
WORKDIR /app

RUN apk add --no-cache build-base pkgconfig

# Copy go.mod and go.sum to leverage Docker's layer caching
COPY go.mod .
COPY go.sum .

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application, creating a static binary
# -o specifies the output file name
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/main ./main.go

# Stage 2: Create a minimal runtime image
FROM alpine:3.19

# Install ca-certificates for HTTPS communication if needed
RUN apk add --no-cache ca-certificates

# Set the working directory for the final application
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main /app/main

# Define the command to run the application when the container starts
ENTRYPOINT ["/app/main"]
CMD ["-config", "/app/config.json"]