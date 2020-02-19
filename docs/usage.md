## Project

## Repository

## Recipe

### Config

A recipe config file is made of two parts:

* a describing manifest (description, files to sync,...) handled by a fixed `manala` map key
* some custom variables serving two purposes:
    * provide default values
    * scaffold validation schema

```yaml
# Manifest
manala:
    description: Saucerful of secrets # Mandatory description
    sync:
      - .manala                       # ".manala" path will be synchronized on project

# Custom variables
foo: bar    # Provide default value for "foo"
bar:        # Scaffold "bar" validation schema as an object
    baz: [] # Scaffold "bar.baz" validation schema as an array
```

### Validation

As seen before, a validation schema is scaffolded from custom variables provided in recipe config file, using [JSON Schema](https://json-schema.org/).

```yaml
bar:
    baz: []
```

generate:

```json
{
  "type": "object",
  "properties": {
    "bar": {
      "type": "object",
      "properties": {
        "baz": {
          "type": "array"
        }
      }
    }
  }
}
```

!!! warning
    Only objects and arrays are supported

Custom validation schema could be provided using doc annotation 

```yaml
bar:
    baz: []
    # @schema {"enum": [null, 123, "foo"]}
    qux: 123
```

```json
{
  "type": "object",
  "properties": {
    "bar": {
      "type": "object",
      "properties": {
        "baz": {
          "type": "array"
        },
        "qux": {
          "enum": [null, 123, "foo"]
        }
      }
    }
  }
}
```

Some custom formats are also provided for the win:

* `go-repo`
* `file-path`
* `domain` 

### Options

Recipe options could be provided using doc annotation. They will be prompted to user during a project initialization.

```yaml
bar:
    baz: []
    # @schema {"enum": [null, 123, "foo"]}
    # @option {"label": "Qux value"} 
    qux: 123
```

Options fields type are guessed by schema details. For instance, an `enum` will  generate a drop down select, and a 
string `type` will generate a text input.

### Content

Recipes support three kind of files:

**Regular**

**Template**

Functions are supplied by the built-in [Go text/template package](https://golang.org/pkg/text/template/) and the
[Sprig template function library](http://masterminds.github.io/sprig/).

Additionally, following functions are provided:
* `toYaml`: serialize variables as yaml

**Dist**
