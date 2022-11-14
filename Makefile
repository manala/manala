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
		go build -o /go/bin/manala
	docker compose run --rm \
		vhs \
		docs/demo/demo.tape
.PHONY: docs/demo

go:
	docker compose run --rm \
		go

mkdocs:
	docker compose run --rm --service-ports \
		mkdocs
