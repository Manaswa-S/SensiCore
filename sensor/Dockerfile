FROM alpine:latest

WORKDIR /root/

COPY sensor .
COPY .example.env .env
COPY start.sh /start.sh

RUN chmod +x /start.sh

ENTRYPOINT [ "./sensor" ]