#!/bin/bash

# Cross compile app
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o teams . || { echo "Building Go App failed" ; exit 1; }

# Build the docker container
docker build -t centralmarinsoccer/teams . || { echo "Building Docker container failed" ; exit 1; }

# Push to docker hub
docker push centralmarinsoccer/teams || { echo "Pushing container to Docker hub failed" ; exit 1; }
