# Use the official Golang image as the base image
FROM golang:1.21.1-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Enable CGO and build the Go application
RUN CGO_ENABLED=1 go build -o main .

# Use a minimal Alpine Linux image for the final stage
FROM alpine:latest

# Install Bash and SQLite
RUN apk add --no-cache bash sqlite

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled binary from the build stage
COPY --from=build /app/main .

# Copy the config.yml file from the build stage
COPY --from=build /app/config.yml .

# Set the entrypoint to run the compiled binary directly
ENTRYPOINT ["./main"]