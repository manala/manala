.SILENT:

include .make/help.mk
include .make/text.mk

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

##########
# MkDocs #
##########

## MkDocs - Generate
mkdocs:
	docker compose run --rm --service-ports \
		mkdocs

## MkDocs - Open mkdocs shell
mkdocs.sh:
	docker compose run --rm --service-ports \
		--entrypoint /bin/ash \
		mkdocs

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

#######
# Web #
#######

## Web - Start web server (PORT)
web: PORT = 9400
web:
	echo Start web server...
	docker compose run --rm \
		--publish $(PORT):$(PORT) \
		go \
		go run \
			-tags web_app_build \
			. \
			web --debug --port $(PORT)
.PHONY: web
