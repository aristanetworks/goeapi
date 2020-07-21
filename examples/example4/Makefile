#!/usr/bin/make
# Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.

GOOS ?= linux
GOARCH ?= amd64

NAME := example4

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(NAME)

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o $(NAME)

clean:
	rm -f $(NAME)

.PHONY: build build-mac clean
