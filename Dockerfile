FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
ADD teams /
VOLUME ["/data-cache"]
EXPOSE 8080
CMD ["/teams"]
