PORT := ${or ${PORT},${PORT},8000}
SSH_AUTH_SOCK := /run/host-services/ssh-auth.sock

## Run local server (on port 8000 by default)
run:
	# https://github.com/squidfunk/mkdocs-material
	docker run --rm -v $(CURDIR):/docs -p ${PORT}:${PORT} squidfunk/mkdocs-material serve -a 0.0.0.0:${PORT}

## Deploy as gh-pages
deploy:
	docker run --rm -it \
		-v ${SSH_AUTH_SOCK}:${SSH_AUTH_SOCK} \
		-e SSH_AUTH_SOCK=${SSH_AUTH_SOCK} \
		-v ${PWD}:/docs \
		squidfunk/mkdocs-material gh-deploy
