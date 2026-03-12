.SILENT:

include .make/help.mk
include .make/text.mk

########
# Lint #
########

## Lint - Lint
lint:
	docker compose run --rm \
		golangci-lint \
		golangci-lint run --verbose

#################
# Golangci-lint #
#################

## Golangci-lint - Open shell
golangci-lint.sh:
	docker compose run --rm \
		golangci-lint \
		/bin/bash

#########
# Build #
#########

## Build - build
build:
	goreleaser build --single-target --auto-snapshot --clean

###########
# Release #
###########

## Release - release
release:
	goreleaser release --snapshot --clean

########
# Docs #
########

## Docs - Generate all
docs: docs/commands docs/demo
.PHONY: docs

## Docs - Generate commands
docs/commands:
	echo Generate docs commands...
	docker compose run --rm \
		go \
		go run . docs docs/commands
.PHONY: docs/commands

## Docs - Generate demo
docs/demo:
	echo Generate docs demo...
	docker compose run --rm \
		go \
		go build -a -v -o /go/bin/manala
	docker compose run --rm \
		vhs \
		docs/demo/demo.tape
.PHONY: docs/demo

############
# Zensical #
############

## Zensical - Generate
zensical:
	docker compose run --rm --service-ports \
		zensical

######
# Go #
######

## Go - Open go shell
go.sh:
	docker compose run --rm \
		go \
		/bin/bash

#######
# Vhs #
#######

## Vhs - Open vhs shell
vhs.sh:
	docker compose run --rm \
		--entrypoint /bin/bash \
		vhs
