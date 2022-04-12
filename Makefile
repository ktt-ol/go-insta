.PHONY: build run install

help:
	@echo To start local: go run ./cmd/insta/insta.go -help
	@echo
	@echo To build for Raspberry Pi:
	@echo - make build
	@echo - make install \(speeds up build, req. writeable GOROOT\)

build:
	GOARM=6 GOOS=linux GOARCH=arm go build -v ./cmd/insta
	mv insta insta-arm-linux

install:
	GOARM=6 GOOS=linux GOARCH=arm go install -v ./cmd/insta