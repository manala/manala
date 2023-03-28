.SILENT:

-include .make/help.mk
-include .make/text.mk

.PHONY: build

## Setup - Install dependencies
install:
	npm install

install@integration:
	npm install --color=always --no-progress --no-audit --no-fund

## Setup - Update dependencies
update:
	npm update

## Serve - Serve
serve:
	npx webpack serve \
		--mode=development \

## Build - Build
build: export BUILD = 1
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
