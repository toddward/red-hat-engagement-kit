# engagement_selection.tape - Select an engagement and verify sidebar updates
Require bash

Output tui/tests/vhs/output/engagement_selection.gif

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
Down
Enter
Sleep 2s
Screenshot tui/tests/vhs/screenshots/engagement_list.png

Enter
Sleep 2s
Screenshot tui/tests/vhs/screenshots/engagement_selected.png

Sleep 1s
Type "q"
