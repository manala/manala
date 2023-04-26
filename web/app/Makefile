.SILENT:

include ../../.make/help.mk
include ../../.make/text.mk

.PHONY: build

## Setup - Install dependencies
install:
	npm install

install@integration:
	npm clean-install --color=always --no-progress --no-audit --no-fund

## Setup - Update dependencies
update:
	npm update

## Serve - Serve (API_PORT, API_BASE_URI)
serve: export API_PORT = 9400
serve: export API_BASE_URI = http://localhost:$(API_PORT)/api
serve:
	npx webpack serve \
		--mode=development

## Build - Build (API_BASE_URI)
build: export BUILD = 1
build: export API_BASE_URI = http://localhost:9400/api
build:
	npx webpack \
		--mode=production

## Tests - Lint
lint: lint.eslint lint.tsc

lint.eslint:
	npx eslint src/* --ext .js,.json,.tsx,.ts --fix

lint.tsc:
	npx tsc --noemit

lint@integration: lint.eslint@integration lint.tsc@integration

lint.eslint@integration:
	npx eslint src/* --ext .js,.json,.tsx,.ts

lint.tsc@integration:
	npx tsc --noemit
