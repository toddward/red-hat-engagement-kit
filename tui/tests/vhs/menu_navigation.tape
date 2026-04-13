# menu_navigation.tape - Navigate through all top-level menu items
Require bash

Output tui/tests/vhs/output/menu_navigation.gif

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

Down
Sleep 500ms
Screenshot tui/tests/vhs/screenshots/menu_item_2.png

Down
Sleep 500ms
Screenshot tui/tests/vhs/screenshots/menu_item_3.png

Down
Sleep 500ms

Down
Sleep 500ms

Up
Up
Up
Up
Sleep 500ms
Screenshot tui/tests/vhs/screenshots/menu_item_1.png

Enter
Sleep 2s
Screenshot tui/tests/vhs/screenshots/skills_submenu.png

Escape
Sleep 500ms
Screenshot tui/tests/vhs/screenshots/back_to_main.png

Sleep 1s
Type "q"
