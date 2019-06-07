FROM golang:alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /go/src/github.com/cybersectech-org/urlinsane
COPY . .

RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/urlinsane ./cmd

FROM iron/go

COPY --from=builder /go/bin/urlinsane /bin/urlinsane

EXPOSE 8080

CMD ["urlinsane", "server", "-a", "0.0.0.0", "-p", "8080", "-s"]
