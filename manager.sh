#!/bin/bash

# Set the name of the Docker image and container
IMAGE_NAME="your-image-name"
CONTAINER_NAME="your-container-name"

# Set the database file path relative to the current working directory
DB_FILE_PATH="$(pwd)/database.db"

# Stop and remove the existing container (if it exists)
if docker ps -a --format '{{.Names}}' | grep -q $CONTAINER_NAME; then
    echo "Stopping and removing the existing container..."
    docker stop $CONTAINER_NAME
    docker rm $CONTAINER_NAME
fi

# Remove the existing image (if it exists)
if docker images --format '{{.Repository}}' | grep -q $IMAGE_NAME; then
    echo "Removing the existing image..."
    docker rmi $IMAGE_NAME
fi

# Build the Docker image
echo "Building the Docker image..."
docker build -t $IMAGE_NAME .

# Run the Docker container with the mounted database file
echo "Running the Docker container..."
docker run --name $CONTAINER_NAME -d -v "$DB_FILE_PATH:/app/database.db" $IMAGE_NAME

echo "Docker container is up and running with the mounted database file!"