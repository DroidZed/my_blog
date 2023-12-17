FROM golang as builder

LABEL authors="DroidZed"

WORKDIR /usr/src/app

EXPOSE 8000

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download && go mod verify

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -v -o bin/golance cmd/go_lance/main.go

# Run the web service on container startup.
CMD ["/usr/src/app/bin/golance"]

# TODO: READ THIS: https://www.ardanlabs.com/blog/2020/02/docker-images-part1-reducing-image-size.html