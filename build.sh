#!/bin/bash

build() {
  echo "Building for Linux..."
  GOOS=linux GOARCH=amd64 go build -o ./bin/docker-cron .
  GOOS=linux GOARCH=386 go build -o ./bin/docker-cron.linux.386 .
  GOOS=linux GOARCH=arm go build -o ./bin/docker-cron.linux.arm .

  echo "Building for Mac..."
  GOOS=darwin GOARCH=amd64 go build -o ./bin/docker-cron.darwin.amd64 .
  GOOS=darwin GOARCH=386 go build -o ./bin/docker-cron.darwin.386 .
}

build
