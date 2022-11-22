FROM alpine:latest

RUN mkdir /app

COPY ./mail.bin /app
COPY templates /templates

CMD ["/app/mail.bin"]
