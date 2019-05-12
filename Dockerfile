
# Start from golang v1.11 base image
FROM golang:1.11

# Add Maintainer Info
LABEL maintainer="Tal Hachi <rangertaha@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/rangertaha/urlinsane

# Copy everything from the current directory
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

RUN mv $GOPATH/bin/cmd $GOPATH/bin/urlinsane

# Exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["urlinsane", "server", "-a", "0.0.0.0", "-p", "8080"]
