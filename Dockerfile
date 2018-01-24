#FROM scratch
#FROM gliderlabs/alpine:3.6
FROM debian
#ADD ca-certificates.crt /etc/ssl/certs/
ADD PineappleScope-linux-amd64 /bin/app
ADD resources/ /resources/

#RUN apk update \
#    && apk add sqlite
#    && apk add bash
RUN apt-get update && apt-get install -y sqlite

EXPOSE 8177

VOLUME /var/db/

ENTRYPOINT ["/resources/entrypoint.sh"]
