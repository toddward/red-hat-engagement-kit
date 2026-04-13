# startup.tape - Verify TUI launches and renders initial state
Require bash

Output tui/tests/vhs/output/startup.gif

Set Shell "bash"
Set Width 120
Set Height 35
Set FontSize 14

Hide
Type "cd $PWD"
Enter
Sleep 500ms
Show

Type "./tui/tui.sh"
Enter

Sleep 3s

Screenshot tui/tests/vhs/screenshots/startup_menu.png

Sleep 1s
Type "q"
