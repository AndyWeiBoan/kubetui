package app

import (
	"os"
	"os/signal"
	"syscall"
	"unsafe"

	"kubetui/tty"
	"kubetui/view"
)

var (
	terminateSigChan chan os.Signal
	windowResizeChan chan os.Signal
	ttyHandler       *tty.TTY
	winSize          *view.Winsize
	rootView         *view.Root
)

func init() {
	terminateSigChan = make(chan os.Signal, 1)
	windowResizeChan = make(chan os.Signal, 1)
	signal.Notify(terminateSigChan, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(windowResizeChan, syscall.SIGWINCH)
	winSize = getWinSize()
	ttyHandler = tty.NewTTY()
	rootView = &view.Root{Name: "root", X: 1, Y: 1, H: 100, W: 100}
	screen := rootView.GetScreen(winSize)
	_ = screen
	ttyHandler.Draw(screen)
	go onWindowsResize()
	go onShutdown()
}

func Start() {
	for byte := range ttyHandler.Read() {
		_ = byte
	}
}

func onShutdown() {
	<-terminateSigChan
	ttyHandler.Close()
	os.Exit(0)
}

func onWindowsResize() {
	for {
		<-windowResizeChan
		winSize = getWinSize()
		screen := rootView.GetScreen(winSize)
		ttyHandler.Draw(screen)
	}
}

// Source - https://stackoverflow.com/a/16576712
func getWinSize() *view.Winsize {
	ws := &view.Winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return ws
}
