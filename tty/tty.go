package tty

import (
	"bufio"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

type bufWriter struct{ *bufio.Writer }

var (
	bufHandler  *bufWriter
	ttyHandler  *os.File
	sig         chan os.Signal
	origTermios syscall.Termios
)

// Source - https://stackoverflow.com/a/16576712
func getWinSize() *winsize {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	log.Printf("row=%v", ws.Row)
	log.Printf("col=%v", ws.Col)
	log.Printf("xpixel=%v", ws.Xpixel)
	log.Printf("ypixel=%v", ws.Ypixel)
	return ws
}

func SetUp() {
	ws := getWinSize()
	// registry term notify
	sig = make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// open tty
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	ttyHandler = tty
	if err != nil {
		panic(err)
	}
	if bufHandler == nil {
		bufHandler = &bufWriter{bufio.NewWriter(tty)}
	}
	setRawMode()
	bufHandler.execCmds(ANSI.Screen.EnableAlternativeBuffer)
	bufHandler.drawPanel(ws)

	// watting for term signal, then return 0
	go handleShoutdown()
}

func Write(s string) {
	bufHandler.WriteString(s)
	bufHandler.Flush()
}

func Close() {
	restoreMode()
	bufHandler.execCmds(ANSI.Screen.DisableAlternativeBuffer)
	ttyHandler.Close()
}

func (w *bufWriter) drawLine(ws *winsize) {
	for i := 1; i < int(ws.Col); i++ {
		// w.WriteString("\033[1;" + strconv.Itoa(i) + "H" + string(0x2501))
	}
	w.WriteString("\033[1;1H" + strings.Repeat("─", int(ws.Col)-1))

	w.Flush()
}

func (w *bufWriter) drawPanel(ws *winsize) {
	// draw upper left corner ┏
	w.WriteString(concatCommands(esc, getPosition(1, 1), string(0x250F)))

	// draw upper right corner ┓
	w.WriteString(concatCommands(esc, getPosition(1, int(ws.Col)), string(0x2513)))

	// draw lower left corner ┗
	w.WriteString(concatCommands(esc, getPosition(int(ws.Row), 1), string(0x2517)))

	// draw lower right corner ┛
	w.WriteString(concatCommands(esc, getPosition(int(ws.Row), int(ws.Col)), string(0x251B)))

	// draw upper and lower horizontal line ━
	horizontalLine := strings.Repeat(string(0x2501), int(ws.Col)-2)
	w.WriteString(concatCommands(esc, getPosition(1, 2), horizontalLine))
	w.WriteString(concatCommands(esc, getPosition(int(ws.Row), 2), horizontalLine))

	// draw left vertical line ┃
	verticalLine := strings.Repeat(string(0x2503)+"\n", int(ws.Row)-2)
	w.WriteString(concatCommands(esc, getPosition(2, 1), verticalLine))

	// drawal right vertical line ┃
	for i := 2; i < int(ws.Row); i++ {
		w.WriteString(concatCommands(esc, getPosition(i, int(ws.Col)), string(0x2503)))
	}

	w.Flush()
}

func getPosition(row int, col int) string {
	return "[" + strconv.Itoa(row) + ";" + strconv.Itoa(col) + "H"
}

func (w *bufWriter) execCmds(cmds ...string) {
	for _, cmd := range cmds {
		w.WriteString(cmd)
	}
	w.Flush()
}

func handleShoutdown() {
	<-sig
	Close()
	os.Exit(0)
}

func setRawMode() {
	// 讀取目前設定
	syscall.Syscall(syscall.SYS_IOCTL,
		ttyHandler.Fd(),
		syscall.TIOCGETA, // 或 TIOCGETA on macOS
		uintptr(unsafe.Pointer(&origTermios)),
	)
	raw := origTermios
	// 關掉 echo 和 canonical mode
	raw.Lflag &^= syscall.ECHO | syscall.ICANON
	// 套用
	syscall.Syscall(syscall.SYS_IOCTL,
		ttyHandler.Fd(),
		syscall.TIOCSETA, // 或 TIOCSETA on macOS
		uintptr(unsafe.Pointer(&raw)),
	)
}

func restoreMode() {
	syscall.Syscall(syscall.SYS_IOCTL,
		ttyHandler.Fd(),
		syscall.TIOCSETA, // 或 TIOCGETA on macOS
		uintptr(unsafe.Pointer(&origTermios)),
	)
}
