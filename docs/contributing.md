## Build

Requirements

* Go 1.13+

Update modules

```shell
go get -u ./...
go mod tidy
```

## Documentation

```shell
docker run --rm -v $(pwd):/data -p 8000:8000 nicksantamaria/mkdocs serve -a 0.0.0.0:8000
```
