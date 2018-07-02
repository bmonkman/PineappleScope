#FROM scratch
# FROM gliderlabs/alpine:3.6
FROM debian

ADD PineappleScope-linux-amd64 /app
ADD resources/ /resources/

ENV TZ=America/Vancouver
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean

EXPOSE 1111

ENV DBFILE /var/db/pineapplescope.db
VOLUME /var/db/

#CMD ["sleep","1000"]
ENTRYPOINT ["/app"]
