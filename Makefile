.SILENT:
.PHONY: docs

## Run documentation local server (on port 8000 by default)
## See: https://github.com/squidfunk/mkdocs-material
docs:
	docker run \
		--rm \
		--volume $(CURDIR):/docs \
		--publish 8000:${or ${PORT},${PORT},8000} \
		squidfunk/mkdocs-material:8.3.8 \
		serve -a 0.0.0.0:8000
