# Use smallest docker image
FROM scratch

ADD ca-certificates.crt /etc/ssl/certs/
ADD static /var/www/api/teams/static

EXPOSE 8080

# Expose the static files
VOLUME ["/var/www/api/teams/static"]

ADD teams /
CMD ["/teams"]
