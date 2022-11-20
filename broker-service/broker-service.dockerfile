FROM alpine:latest

RUN mkdir /app

COPY ./brokerApp.bin /app/brokerApp.bin

CMD ["/app/brokerApp.bin"]
