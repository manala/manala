.SILENT:
.PHONY: docs

## Run documentation local server (on port 8000 by default)
## See: https://github.com/squidfunk/mkdocs-material
docs:
	docker run \
		--rm \
		--volume $(CURDIR):/docs \
		--publish 8000:${or ${PORT},${PORT},8000} \
		squidfunk/mkdocs-material:5.5.12 \
		serve -a 0.0.0.0:8000
