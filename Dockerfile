FROM golang
Add . /go/src/github.com/centralmarinsoccer/teams
WORKDIR /go/src/github.com/centralmarinsoccer/teams
RUN go get ./...
VOLUME ["/go/src/github.com/centralmarinsoccer/teams/data-cache"]

ENTRYPOINT /go/bin/teams
EXPOSE 8080
