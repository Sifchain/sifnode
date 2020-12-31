# Makefile

# Speed up re-running the tests by only doing a go build if any .go file has changed

SHELL:=/bin/bash

gofiles=$(shell find ${BASEDIR}/* -name \*.go) go.mod go.sum

.PHONY: $(gofiles)
vagrant/.goBuild: $(gofiles)
	make -C ${BASEDIR} install
	touch $@
