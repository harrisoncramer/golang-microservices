FROM alpine:latest

RUN mkdir /app

COPY ./mail.bin /app

CMD ["/app/mail.bin"]
