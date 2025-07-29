FROM alpine:latest

WORKDIR /root/

COPY server .
COPY .example.env .env
COPY start.sh /start.sh

RUN chmod +x /start.sh

EXPOSE 8686

ENTRYPOINT [ "./server" ]