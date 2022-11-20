FROM alpine:latest

RUN mkdir /app

COPY ./loggerApp.bin /app

CMD ["/app/loggerApp.bin"]
