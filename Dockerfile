# Use the latest version of the official Golang image as the base image
FROM golang:1.20

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files to the working directory
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the application source code to the working directory
COPY . .

# Build the Golang application
RUN GOARCH=arm64 go build -o main .

# Mount the current directory as a volume to persist file changes
VOLUME ["/app"]

# Expose the port on which your application will run (replace 8080 with your desired port)
EXPOSE 8080
RUN chmod +x main

# Set the command to run your application
CMD ["./main"]