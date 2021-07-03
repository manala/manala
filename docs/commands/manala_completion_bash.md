## manala completion bash

generate the autocompletion script for bash

### Synopsis


Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:
$ source <(manala completion bash)

To load completions for every new session, execute once:
Linux:
  $ manala completion bash > /etc/bash_completion.d/manala
MacOS:
  $ manala completion bash > /usr/local/etc/bash_completion.d/manala

You will need to start a new shell for this setup to take effect.
  

```
manala completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -c, --cache-dir string   use cache directory
  -d, --debug              set debug mode
```

### SEE ALSO

* [manala completion](manala_completion.md)	 - generate the autocompletion script for the specified shell

