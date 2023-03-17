.SILENT:

docs: docs/commands docs/demo
.PHONY: docs

docs/commands:
	echo Generate docs commands...
	docker compose run --rm \
		go \
		go run . docs docs/commands
.PHONY: docs/commands

docs/demo:
	echo Generate docs demo...
	docker compose run --rm \
		go \
		go build -a -v -o /go/bin/manala
	docker compose run --rm \
		vhs \
		docs/demo/demo.tape
.PHONY: docs/demo

go.sh:
	docker compose run --rm \
		go \
		/bin/bash

mkdocs:
	docker compose run --rm --service-ports \
		mkdocs

mkdocs.sh:
	docker compose run --rm --service-ports \
		--entrypoint /bin/ash \
		mkdocs

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
