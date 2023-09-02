# Use an official Go runtime as a base image
FROM golang:1.18-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application code to the container's workspace
COPY . .
RUN CGO_ENABLED=0
RUN apk add --no-cache libc6-compat 

# Build the Go application
RUN go build -o fmtbe ./cmd/FindMeTime/

# Expose a port that the application will listen on
EXPOSE 8080

# Command to run the executable
# CMD ["ls"]
ENTRYPOINT ["./fmtbe"]
