# Start from golang base image
FROM golang:1.17.6-alpine3.14 as BenderBeer

# Add Maintainer info
LABEL maintainer="Yael Castro <contacto@yael-castro.com>"

WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod tidy

# Build the Go app
RUN env GOOS=linux go build main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]