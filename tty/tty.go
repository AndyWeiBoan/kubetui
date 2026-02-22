package tty

import (
	"bufio"
	"iter"
	"os"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

type bufWriter struct{ *bufio.Writer }

type TTY struct {
	bufWriter   *bufWriter
	file        *os.File
	origTermios syscall.Termios
}

type Cell struct {
	Char    rune
	Changed bool
}

type Screen struct {
	Cells [][]Cell
}

func NewTTY() *TTY {
	file, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}

	tty := &TTY{
		file:      file,
		bufWriter: &bufWriter{bufio.NewWriter(file)},
	}
	tty.makeRawMode()
	tty.Exec(
		ANSI.Screen.EnableAlternativeBuffer,
		ANSI.Keyboard.EnableKeyboardEvent,
	)
	return tty
}

func (t *TTY) Close() {
	t.restoreOriginMode()
	t.Exec(ANSI.Screen.DisableAlternativeBuffer)
	t.file.Close()
}

func (t *TTY) Draw(screen Screen) {
	var cmd string
	for i := 1; i < len(screen.Cells); i++ {
		for j := 1; j < len(screen.Cells[i]); j++ {
			cell := screen.Cells[i][j]
			if !cell.Changed {
				continue
			}
			runeCmd := concatCommands(esc, getPosition(i, j), string(cell.Char))
			cmd = concatCommands(cmd, runeCmd)
		}
	}
	t.Exec(cmd)
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

func getPosition(row int, col int) string {
	return "[" + strconv.Itoa(row) + ";" + strconv.Itoa(col) + "H"
}

func (t *TTY) Exec(cmds ...string) {
	for _, cmd := range cmds {
		t.bufWriter.WriteString(cmd)
	}
	t.bufWriter.Flush()
}

func (t *TTY) makeRawMode() {
	// 讀取目前設定
	syscall.Syscall(syscall.SYS_IOCTL,
		t.file.Fd(),
		syscall.TIOCGETA, // 或 TIOCGETA on macOS
		uintptr(unsafe.Pointer(&t.origTermios)),
	)
	raw := t.origTermios
	// 關掉 echo 和 canonical mode
	raw.Lflag &^= syscall.ECHO | syscall.ICANON
	// 套用
	syscall.Syscall(syscall.SYS_IOCTL,
		t.file.Fd(),
		syscall.TIOCSETA, // 或 TIOCSETA on macOS
		uintptr(unsafe.Pointer(&raw)),
	)
}

func (t *TTY) restoreOriginMode() {
	syscall.Syscall(syscall.SYS_IOCTL,
		t.file.Fd(),
		syscall.TIOCSETA, // 或 TIOCGETA on macOS
		uintptr(unsafe.Pointer(&t.origTermios)),
	)
}

func (t *TTY) Read() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for {
			buffer := make([]byte, 256)
			length, err := t.file.Read(buffer)
			if err != nil {
				continue
			}
			data := buffer[:length]
			if length == 0 {
				time.Sleep(1 * time.Second)
				continue
			}
			if !yield(data) {
				return
			}
		}
	}
}
