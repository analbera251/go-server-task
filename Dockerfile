# Use the official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod file to the container
COPY go.mod .

# Download dependencies
RUN go mod download

# Copy the local code to the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port on which the Go application will run
EXPOSE 8000

# Command to run the executable
CMD ["./main"]
