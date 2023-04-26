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

## Web - Install
web.install:
	go install
	$(MAKE) --directory web/app install
.PHONY: install

web: web.serve

## Web - Start api web server + React app (PORT)
web.serve: PORT = 9400
web.serve:
	# https://www.npmjs.com/package/concurrently
	npx concurrently "make web.serve.api+local PORT=$(PORT)" "make web.serve.front API_PORT=$(PORT)" \
		--names="API,Front" \
		--prefix=name \
		--prefix-colors="auto" \
 		--kill-others \
 		--kill-others-on-fail \
 		--colors
.PHONY: web.serve

## Web - Start front (React) app only (API_PORT)
web.serve.front: API_PORT = 9400
web.serve.front:
	echo Start front server...
	$(MAKE) --directory web/app serve API_PORT=$(API_PORT)
.PHONY: web.serve.front

## Web - Start api web server only (requires go) (PORT)
web.serve.api+local: PORT = 9400
web.serve.api+local:
	echo Start api server...
	go run \
		-tags web_app_build \
		. \
		web --debug --port $(PORT)
.PHONY: web.serve.api+local

## Web - Start api web server only using Docker (PORT)
web.serve.api+docker: PORT = 9400
web.serve.api+docker:
	echo Start api server...
	docker compose run --rm \
		--publish $(PORT):$(PORT) \
		go \
		go run \
			-tags web_app_build \
			. \
			web --debug --port $(PORT)
.PHONY: web.serve.api+docker
