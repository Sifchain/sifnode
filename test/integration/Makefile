# Makefile

# Speed up re-running the tests by only doing a go build if any .go file has changed

SHELL:=/bin/bash


ifndef BASEDIR
$(error BASEDIR is not set)
endif

ifndef GOBIN
$(error GOBIN is not set)
endif

gofiles=$(shell find ${BASEDIR}/* -name \*.go) ${BASEDIR}/go.mod ${BASEDIR}/go.sum

${GOBIN}/sifnoded: $(gofiles)
	make -C ${BASEDIR} install
	touch $@
