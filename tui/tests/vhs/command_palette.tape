# command_palette.tape - Open palette, type to filter, select item
Require bash

Output tui/tests/vhs/output/command_palette.gif

Set Shell "bash"
Set Width 120
Set Height 35

Hide
Type "cd $PWD"
Enter
Sleep 500ms
Show

Type "./tui/tui.sh"
Enter
Sleep 3s

Type "/"
Sleep 1s
Screenshot tui/tests/vhs/screenshots/palette_open.png

Type "setup"
Sleep 500ms
Screenshot tui/tests/vhs/screenshots/palette_filtered.png

Escape
Sleep 500ms
Screenshot tui/tests/vhs/screenshots/palette_closed.png

Sleep 1s
Type "q"
