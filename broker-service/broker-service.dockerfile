FROM alpine:latest

RUN mkdir /app

COPY ./broker.bin /app/broker.bin

CMD ["/app/broker.bin"]
