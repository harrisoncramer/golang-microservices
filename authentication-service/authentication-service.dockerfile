FROM alpine:latest

RUN mkdir /app

COPY ./authApp.bin /app

CMD ["/app/authApp.bin"]
