# Start from the official Go image
FROM golang:latest

# Set the current working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Install the dependencies
RUN go mod download

# Build the Go application
RUN go build -o app

# Expose port 8080 for the container
EXPOSE 1234

# Start the application when the container starts
CMD ["./app"]
