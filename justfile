#!/usr/bin/env just --justfile

set shell := ["zsh", "-cu"] 

# Lists the justfile commands
@default:
  @just --list

# Build the package
@build:
  go build -ldflags "-s -w"

@run: build
	./extip-svr

# Be a good citizen
@fmt:
  go fmt
