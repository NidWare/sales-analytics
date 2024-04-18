# Use the latest version of the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files to the working directory
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the application source code to the working directory
COPY . .

# Build the Golang application
RUN go build -o main .

# Expose the port on which your application will run (replace 8080 with your desired port)
EXPOSE 8080

# Mount the current directory as a volume to persist file changes
VOLUME [ "/app" ]

# Set the command to run your application
CMD ["./main"]