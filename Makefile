VERSION ?= dev

BINARY_NAME := gix

PKG := github.com/ademajagon/gix/cmd

build:
	go build -ldflags "-X $(PKG).version=$(VERSION)" -o $(BINARY_NAME) .

clean:
	rm -f $(BINARY_NAME)