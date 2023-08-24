## manala completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(manala completion zsh)

To load completions for every new session, execute once:

#### Linux:

	manala completion zsh > "${fpath[1]}/_manala"

#### macOS:

	manala completion zsh > $(brew --prefix)/share/zsh/site-functions/_manala

You will need to start a new shell for this setup to take effect.


```
manala completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -c, --cache-dir string   use cache directory
  -d, --debug              set debug mode
```

### SEE ALSO

* [manala completion](manala_completion.md)	 - Generate the autocompletion script for the specified shell

