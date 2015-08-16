.phony: clean

PROJECT = chromos
SOURCES = $(wildcard src/*.go)
ASSETS  = $(wildcard src/resources/*)
VERSION = $(shell cat VERSION)
REPLACE = s/9\.9\.9/$(VERSION)/g
UNAME   = $(shell uname)

ifeq ($(UNAME), Darwin)
	PLATFORM := osx
else ifeq ($(UNAME), Windows)
	PLATFORM := win
else
	PLATFORM := nix
endif

include Makefile.$(PLATFORM)

clean:
	rm -rf build
