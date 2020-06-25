## manala init

Init project

### Synopsis

Init will init a project.

Example: manala init -> resulting in a project init in a directory (default to the current directory)

```
manala init [dir] [flags]
```

### Options

```
  -h, --help            help for init
  -i, --recipe string   recipe
```

### Options inherited from parent commands

```
  -c, --cache-dir string    cache directory
  -d, --debug               debug mode
  -o, --repository string   repository
```

## manala list

List recipes

### Synopsis

List will list recipes available on repository.

Example: manala list -> resulting in a recipes list display

```
manala list [flags]
```

### Options

```
  -h, --help   help for list
```

### Options inherited from parent commands

```
  -c, --cache-dir string    cache directory
  -d, --debug               debug mode
  -o, --repository string   repository
```

## manala update

Update project

### Synopsis

Update will update project, based on recipe and related variables defined in manala.yaml.

Example: manala update -> resulting in an update in a directory (default to the current directory)

```
manala update [dir] [flags]
```

### Options

```
  -h, --help   help for update
```

### Options inherited from parent commands

```
  -c, --cache-dir string    cache directory
  -d, --debug               debug mode
  -o, --repository string   repository
```

## manala watch

Watch project

### Synopsis

Watch will watch project, and launch update on changes.

Example: manala watch -> resulting in a watch in a directory (default to the current directory)

```
manala watch [dir] [flags]
```

### Options

```
  -h, --help     help for watch
  -n, --notify   use system notifications
  -i, --recipe   watch recipe too
```

### Options inherited from parent commands

```
  -c, --cache-dir string    cache directory
  -d, --debug               debug mode
  -o, --repository string   repository
```
