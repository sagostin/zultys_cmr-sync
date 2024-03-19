# Use an official Go runtime as a parent image
FROM golang:1.21.1-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace.
COPY . /app

# Download all the dependencies
RUN go mod download

# Build the Go app
RUN go build -o main .

# Set environment variables
ENV CONFIG_PATH=./config.json

# Command to run the executable
CMD ["/app/main"]