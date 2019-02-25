# First container builds the app but is rather large
FROM golang:1.11 as builder
LABEL maintainer="Mattias Forsman <screenshooter@mattiasforsman.se>"
WORKDIR $GOPATH/src/github.com/traed/screenshooter

COPY . .

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/scsh .

# Use a separete image to run the app. No need for the Go runtime here.
FROM alpine:latest

RUN apk --no-cache update && apk --no-cache add ca-certificates

# Installs latest Chromium package.
RUN echo @edge http://nl.alpinelinux.org/alpine/edge/community >> /etc/apk/repositories \
    && echo @edge http://nl.alpinelinux.org/alpine/edge/main >> /etc/apk/repositories \
    && apk add --no-cache \
    chromium@edge \
    harfbuzz@edge \
    nss@edge \
    && rm -rf /var/cache/*

WORKDIR /root/
COPY --from=builder /go/bin/scsh .

RUN mkdir uploads

EXPOSE 8080

ENV CHROME_PATH /usr/bin/chromium-browser

VOLUME ["/root/uploads"]
CMD ["./scsh", "-d=./uploads", "-t=15"]