FROM alpine:latest

RUN mkdir /app

COPY ./listener.bin /app

CMD ["/app/listener.bin"]
