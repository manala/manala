## manala update

Update project

### Synopsis

Update (manala update) will update project, based on
recipe and related variables defined in manala.yaml.

Example: manala update -> resulting in an update in a directory (default to the current directory)

```
manala update [dir] [flags]
```

### Options

```
  -h, --help                help for update
  -i, --recipe string       force recipe
  -r, --recursive           recursive
  -o, --repository string   force repository
```

### Options inherited from parent commands

```
  -c, --cache-dir string   cache directory (default "/Users/florian.rey/Library/Caches")
  -d, --debug              debug mode (default true)
```

### SEE ALSO

* [manala](manala.md)	 - Let your project's plumbing up to date

