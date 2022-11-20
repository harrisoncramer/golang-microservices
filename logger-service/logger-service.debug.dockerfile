#################
# Compile Stage #
#################
FROM golang:1.17 AS build-env
RUN go install github.com/go-delve/delve/cmd/dlv@latest

###############
# Debug Stage #
###############
FROM alpine:latest
RUN apk add gcompat # Required for delve to run

COPY --from=build-env /go/bin/dlv /bin

RUN mkdir /app

COPY ./loggerDebug /app/loggerDebug

CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app/loggerDebug"]
