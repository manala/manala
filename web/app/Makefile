.SILENT:

.PHONY: build

## Setup - Install dependencies
install:
	npm install

## Setup - Update dependencies
update:
	npm update

## Serve
serve:
	npx webpack serve \
		--mode=development \

## Build
build: export BUILD = 1
build:
	npx webpack \
		--mode=production
