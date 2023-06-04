FROM golang:alpine AS build

# Set necessary environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -tags timetzdata -o main ./

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Build a small image
FROM scratch
LABEL maintainer="Amy Nagle <kabi-git@openmuffin.com>"
LABEL org.opencontainers.image.authors="Amy Nagle <kabi-git@openmuffin.com>"

# Copy the ca-certs and tz data from build image
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy our binary from the build image
COPY --from=build /dist/main /

# Command to run
ENTRYPOINT ["/main"]
