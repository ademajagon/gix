VERSION ?= dev

BINARY_NAME := toka

PKG := github.com/ademajagon/toka/cmd

build:
	go build -ldflags "-X $(PKG).version=$(VERSION)" -o $(BINARY_NAME) .

clean:
	rm -f $(BINARY_NAME)