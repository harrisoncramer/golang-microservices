FROM alpine:latest

RUN mkdir /app

COPY ./logger.bin /app

CMD ["/app/logger.bin"]
