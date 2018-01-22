FROM scratch
#ADD ca-certificates.crt /etc/ssl/certs/
ADD main /
EXPOSE 8177
CMD ["/main"]
