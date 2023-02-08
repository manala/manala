## manala update

Synchronize project(s)

### Synopsis

Update (manala update) will synchronize project(s), based on
repository's recipe and related variables defined in manifest (.manala.yaml).

Example: manala update -> resulting in an update in a project dir (default to the current directory)

```
manala update [dir] [flags]
```

### Options

```
  -h, --help                help for update
  -i, --recipe string       use recipe
  -r, --recursive           set recursive mode
      --ref string          use repository ref
  -o, --repository string   use repository
```

### Options inherited from parent commands

```
  -c, --cache-dir string   use cache directory
  -d, --debug              set debug mode
```

### SEE ALSO

* [manala](manala.md)	 - Let your project's plumbing up to date

