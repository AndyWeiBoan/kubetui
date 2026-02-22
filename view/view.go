package view

import "kubetui/tty"

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
	for i := 0; i < int(winSize.Col); i++ {
		screen.Cells = make([][]tty.Cell, winSize.Col+1)
	}
	return screen
}
