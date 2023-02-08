.SILENT:

docs: docs/commands docs/demo
.PHONY: docs

docs/commands:
	printf "Generate docs commands...\n"
	docker compose run --rm \
		go \
		go run . docs docs/commands
.PHONY: docs/commands

docs/demo:
	printf "Generate docs demo...\n"
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
