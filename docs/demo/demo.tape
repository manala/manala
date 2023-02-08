Output docs/demo/demo.gif

Set FontSize 16
Set Width 1024
Set Height 768
Set Framerate 20
Set Padding 10
Set TypingSpeed 100ms

# Setup
Hide
Type@1ms "export PATH=/go/bin:$PATH"
Enter
Type@1ms "rm -Rf my-project"
Enter
Type@1ms "apt update && apt install -y git"
Enter
Sleep 6s
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
