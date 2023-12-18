Output docs/demo/demo.gif

Set Width 1024
Set Height 768
Set Padding 10
Set WindowBar Colorful

Set FontSize 16

Set Framerate 20
Set TypingSpeed 100ms

# Setup
Hide
Type@1ms "export PATH=/go/bin:$PATH"
Enter
Type@1ms "rm -Rf my-project"
Enter
Type@1ms "apt update && apt install -y git"
Enter
Sleep 10s
Type@1ms "clear"
Enter
Show

Type "manala"
Sleep 1s
Enter
Sleep 3s

Ctrl+L

Type "manala ls --help"
Sleep 1s
Enter
Sleep 3s

Type "manala ls"
Sleep 1s
Enter
Sleep 3s

Type "manala ls -o examples/repository"
Sleep 1s
Enter
Sleep 3s

Ctrl+L

Type "manala init --help"
Sleep 1s
Enter
Sleep 3s

Type "manala init examples/project -o examples/repository"
Sleep 1s
Enter
Sleep 1s

Down@500ms
Down@500ms
Up@500ms
Enter@1s
Sleep 1s

Type "Sir Galahad"
Tab@1s
Type "I seek the Grail"
Tab@1s
Right@1s
Tab@1s
Down@250ms
Down@250ms
Down@250ms
Down@250ms
Down@250ms
Down@250ms
Down@250ms
Down@500ms
Enter@1s
Enter@1s
Sleep 2s

Type "cd examples/project"
Sleep 1s
Enter

Type "ls -sa"
Sleep 1s
Enter

Type "manala up -o ../repository"
Sleep 1s
Enter

Sleep 4s

# Teardown
Hide
Type@1ms "cd ../.. && rm -Rf examples/project"
Enter
Show
