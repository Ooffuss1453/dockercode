# Use the official GoLang image as the base
FROM golang:1.17.6

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download and cache Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o app

# Expose the desired port
EXPOSE 3000

# Define the command to run the application
CMD ["./app", "-folder", "bookingsystemv2"]

