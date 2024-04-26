# Use a specific version of the golang base image
FROM golang:1.22.1 as builder

# Update CA certificates to ensure HTTPS requests can be made (if your Go code makes any)
RUN update-ca-certificates

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy the source code into the container
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download


# Build the application
RUN go build -o bitcoin-handshake .

# Use a minimal image
FROM scratch

# Import the CA certificates from the builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Set the working directory
WORKDIR /opt

# Copy the pre-built binary file from the previous stage
COPY --from=builder /build/bitcoin-handshake .

# Command to run the executable
CMD ["./bitcoin-handshake"]
