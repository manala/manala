## manala completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	manala completion fish | source

To load completions for every new session, execute once:

	manala completion fish > ~/.config/fish/completions/manala.fish

You will need to start a new shell for this setup to take effect.


```
manala completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -c, --cache-dir string   use cache directory
  -d, --debug              set debug mode
```

### SEE ALSO

* [manala completion](manala_completion.md)	 - Generate the autocompletion script for the specified shell

