FROM alpine:latest

RUN mkdir /app

COPY ./authentication.bin /app

CMD ["/app/authentication.bin"]
