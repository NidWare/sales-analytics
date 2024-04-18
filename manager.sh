#!/bin/bash
# Set the name of your Docker image and container
IMAGE_NAME="your-app-name"
CONTAINER_NAME="your-app-container"

# Function to start the container
start_container() {
    echo "Stopping and removing existing container (if any)..."
    docker stop $CONTAINER_NAME || true
    docker rm $CONTAINER_NAME || true

    echo "Building Docker image..."
    docker build -t $IMAGE_NAME .

    echo "Starting container..."
    docker run -d --name $CONTAINER_NAME -p 8080:8080 -v $(pwd):/app $IMAGE_NAME

    echo "Checking container logs..."
    docker logs $CONTAINER_NAME

    echo "Checking container status..."
    docker ps -a | grep $CONTAINER_NAME
}

# Function to stop the container
stop_container() {
    echo "Stopping container..."
    docker stop $CONTAINER_NAME || true
    echo "Container stopped successfully."
}

# Check the command line argument and call the corresponding function
case "$1" in
    start)
        start_container
        ;;
    stop)
        stop_container
        ;;
    *)
        echo "Usage: $0 {start|stop}"
        exit 1
        ;;
esac