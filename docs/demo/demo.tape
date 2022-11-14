Output docs/demo/demo.gif

Set FontSize 20
Set Width 1600
Set Height 800
Set Framerate 30
Set Padding 12
Set TypingSpeed 100ms

# Setup
Hide
Type@1ms "export PATH=/go/bin:$PATH && rm -Rf my-project && clear"
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

Type "manala ls -o docs/examples/repository"
Sleep 1s
Enter
Sleep 2s

Ctrl+L

Type "manala init --help"
Sleep 1s
Enter
Sleep 3s

Type "manala init my-project -o docs/examples/repository"
Sleep 1s
Enter

Down@1s
Up@1s
Down@1s
Enter@1s
Sleep 2s

Type "my project"
Tab@1s
Down@1s
Down@1s
Enter@1s
Tab@1s
Enter@1s
Sleep 2s

Type "cd my-project"
Sleep 1s
Enter

Type "ls -sa"
Sleep 1s
Enter

Type "manala up -o ../docs/examples/repository"
Sleep 1s
Enter

Sleep 4s

# Teardown
Hide
Type@1ms "cd .. && rm -Rf my-project"
Enter
Show
