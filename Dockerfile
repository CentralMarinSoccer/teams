# Use smallest docker image
FROM scratch

ADD ca-certificates.crt /etc/ssl/certs/
ADD static /static

# Store cache externally so it can persist between container restarts
VOLUME ["/data-cache"]
EXPOSE 8080
# Provide a mechanism to quickly iterate on the client code without having to rebuild the container
VOLUME ["/static"]
ADD teams /
CMD ["/teams"]
