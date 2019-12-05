# Manala

[![Release](https://img.shields.io/github/release/manala/manala.svg?style=flat-square)](https://github.com/manala/manala/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE)
[![Travis](https://img.shields.io/travis/manala/manala.svg?style=flat-square)](https://travis-ci.org/manala/manala)
[![SayThanks.io](https://img.shields.io/badge/SayThanks.io-%E2%98%BC-1EAEDB.svg?style=flat-square)](https://saythanks.io/to/manala)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)

## Install / Update

**homebrew tap:**
```
brew install manala/tap/manala
```

**deb/rpm:**

Download the `.deb` or `.rpm` from the releases page and install with `dpkg -i` and `rpm -i` respectively.

**Shell script:**
```
curl -sfL https://raw.githubusercontent.com/manala/manala/master/godownloader.sh | sh
```

## Build

Requirements

* Go 1.13+

Update modules

```
go get -u ./...
go mod tidy
```

## Usage

### Project

### Repository

### Recipe

Recipes support three kind of files:

**Regular**

**Template**

Functions are supplied by the built-in [Go text/template package](https://golang.org/pkg/text/template/) and the
[Sprig template function library](http://masterminds.github.io/sprig/).

Additionally, following functions are provided:
* `toYaml`: serialize variables as yaml

**Dist**
