ARG GO_VERSION=1.23.1-bullseye

FROM golang:${GO_VERSION} AS build

WORKDIR /build/src

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o my_blog ./cmd/my_blog

FROM scratch

VOLUME /home/my_blog

COPY --from=build /build/src/my_blog /usr/bin/my_blog

EXPOSE 8000

ENTRYPOINT ["/usr/bin/my_blog"]
