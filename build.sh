#!/bin/bash

build() {
  echo "Building for Linux..."
  GOOS=linux GOARCH=amd64 go build -o ./bin/linux/amd64/docker-cron .
  GOOS=linux GOARCH=386 go build -o ./bin/linux/386/docker-cron .
  GOOS=linux GOARCH=arm go build -o ./bin/linux/arm/docker-cron .

  echo "Building for Mac..."
  GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin/amd64/docker-cron .
  GOOS=darwin GOARCH=386 go build -o ./bin/darwin/386/docker-cron .

  echo "Building for FreeBSD..."
  GOOS=freebsd GOARCH=amd64 go build -o ./bin/freebsd/amd64/docker-cron .
  GOOS=freebsd GOARCH=386 go build -o ./bin/freebsd/386/docker-cron .

  echo "Building for Windows..."
  GOOS=windows GOARCH=amd64 go build -o ./bin/windows/amd64/docker-cron .
  GOOS=windows GOARCH=386 go build -o ./bin/windows/386/docker-cron .
}

build
