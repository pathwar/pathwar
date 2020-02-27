# Start from golang:1.14-alpine base image
FROM golang:1.14-alpine

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Set the Current Working Directory inside the container
WORKDIR /go/src/pathwar.land/tool/screenshoter

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# COPY src files
COPY . .

# Build the Go app
RUN go install

# Run app
ENTRYPOINT ["screenshoter"]
