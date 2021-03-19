## Project

## Repository

## Recipe

### Config

A recipe manifest file is made of two parts:

* a config block (description, project manifest template, files to sync,...) handled by a fixed `manala` map key
* some custom variables serving two purposes:
    * provide default values
    * scaffold validation schema

```yaml
# Config
manala:
    description: Saucerful of secrets # Mandatory description
    template: .manala.yaml.tmpl       # Optional project manifest template 
    sync:
      - .manala                       # ".manala" path will be synchronized on project

# Custom variables
foo: bar    # Provide default value for "foo"
bar:        # Scaffold "bar" validation schema as an object
    baz: [] # Scaffold "bar.baz" validation schema as an array
```

### Validation

As seen before, a validation schema is scaffolded from custom variables provided in recipe manifest file, using [JSON Schema](https://json-schema.org/).

```yaml
foo:
    bar: []
baz: 123
qux: {}
```

generate:

```json
{
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "foo": {
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "bar": {
          "type": "array"
        }
      }
    },
    "baz": {},
    "qux": {
      "type": "object",
      "additionalProperties": true,
      "properties": {}
    }
  }
}
```

!!! note
    `additionalProperties` default value depends on the number of the relating object properties. An empty object
    (zero properties) will lead to a default `true`, meaning that all properties are left to the discretion of the end
    user. Conversely, one or more properties will lead to a default `false`.

!!! warning
    Only objects and arrays are auto scaffolded

Custom validation schema could be provided using doc annotation 

```yaml
foo:
    bar: []
# @schema {"enum": [null, 123, "foo"]}
baz: 123
# @schema {"additionalProperties": false}
qux: {}
```

```json
{
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "foo": {
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "bar": {
          "type": "array"
        }
      }
    },
    "baz": {
      "enum": [null, 123, "foo"]
    },
    "qux": {
      "type": "object",
      "additionalProperties": false,
      "properties": {}
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

In case of an `enum`, choices ares available from left to right, first one will be default.

### Content

Recipes support three kind of files:

**Regular**

**Template**

Functions are supplied by the built-in [Go text/template package](https://golang.org/pkg/text/template/) and the
[Sprig template function library](http://masterminds.github.io/sprig/).

Additionally, following functions are provided:
* `toYaml`: serialize variables as yaml

**Dist**
