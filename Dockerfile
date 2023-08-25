FROM golang:1.20-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the Working Directory inside the container
COPY main.go .

# Install necessary packages and dependencies
RUN apk add --no-cache git
RUN go mod init aws-rds-cert-checker
RUN go get github.com/aws/aws-sdk-go

# Build the Go app
RUN go build -o aws-rds-cert-checker .

# Run the binary program produced by `go build`
ENTRYPOINT ["./aws-rds-cert-checker"]

