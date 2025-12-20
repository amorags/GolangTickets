# Start from the official Go image
# Using 1.25-alpine based on your go.mod
FROM golang:1.25-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker cache
COPY go.mod ./
# COPY go.sum ./ 
# Note: go.sum is commented out until you add dependencies

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
# We are building the package located in cmd/api
RUN go build -o main ./cmd/api

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
