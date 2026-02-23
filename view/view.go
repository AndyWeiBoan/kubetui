package view

import (
	"kubetui/tty"
)

type Winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}
type Root struct {
	X    int     // x-coordinate
	Y    int     // y-coordinate
	H    float64 // height / winSize.Row
	W    float64 // width / winSize.Col
	Name string
	// Children []RootPanel
	// Data     []string
}

func (r *Root) GetScreen(winSize *Winsize) tty.Screen {
	screen := tty.Screen{Cells: make([][]tty.Cell, winSize.Row+1)}
	for i := 0; i < int(winSize.Row+1); i++ {
		screen.Cells[i] = make([]tty.Cell, winSize.Col+1)
	}
	for i := 2; i < len(screen.Cells[1]); i++{//  range screen.Cells[1] {
		screen.Cells[1][i].Changed = true
		screen.Cells[1][i].Char = rune(0x2501)
		screen.Cells[len(screen.Cells) -1][i].Changed = true
		screen.Cells[len(screen.Cells) -1][i].Char = rune(0x2501)
	}
	for i := 2;i< len(screen.Cells); i++ {//range screen.Cells {
		screen.Cells[i][1].Changed = true
		screen.Cells[i][1].Char = rune(0x2503)
		screen.Cells[i][len(screen.Cells[i]) -1].Changed = true
		screen.Cells[i][len(screen.Cells[i]) -1].Char = rune(0x2503)
	}
	screen.Cells[1][1].Changed = true
	screen.Cells[1][1].Char = rune(0x250F)
	screen.Cells[1][len(screen.Cells[1])-1].Changed = true
	screen.Cells[1][len(screen.Cells[1])-1].Char = rune(0x2513)
	screen.Cells[len(screen.Cells)-1][1].Changed = true
	screen.Cells[len(screen.Cells)-1][1].Char = rune(0x2517)

	lastRow := screen.Cells[len(screen.Cells)-1]
	lastRow[len(lastRow)-1].Changed = true
	lastRow[len(lastRow)-1].Char = rune(0x251B)
	return screen
}


// func (t *TTY) drawPanel(ws *Winsize) {
// 	horizontalLine := strings.Repeat(string(0x2501), int(ws.Col)-2)
// 	verticalLine := strings.Repeat(string(0x2503)+"\n", int(ws.Row)-2)
//
// 	cmd := concatCommands(
// 		// draw upper left corner ┏
// 		concatCommands(esc, getPosition(1, 1), string(0x250F)),
//
// 		// draw upper right corner ┓
// 		concatCommands(esc, getPosition(1, int(ws.Col)), string(0x2513)),
//
// 		// draw lower left corner ┗
// 		concatCommands(esc, getPosition(int(ws.Row), 1), string(0x2517)),
//
// 		// draw lower right corner ┛
// 		concatCommands(esc, getPosition(int(ws.Row), int(ws.Col)), string(0x251B)),
//
// 		// draw upper and lower horizontal line ━
// 		concatCommands(esc, getPosition(1, 2), horizontalLine),
// 		concatCommands(esc, getPosition(int(ws.Row), 2), horizontalLine),
// 		concatCommands(esc, getPosition(2, 1), verticalLine),
// 	)
// 	for i := 2; i < int(ws.Row); i++ {
// 		// drawal right vertical line ┃
// 		cmd = concatCommands(cmd, concatCommands(esc, getPosition(i, int(ws.Col)), string(0x2503)))
// 	}
// 	t.exec(cmd)
// }


