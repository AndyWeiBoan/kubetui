package main

// "bufio"
// "os"
// "os/signal"
// "syscall"
// "time"

import (
	"fmt"
	"kubetui/app"
)

func main2() {
	fmt.Println(string(1))
}

func main() {
	app.Start()
	// ttyHandler, err := tty.NewTTY()
	// if err != nil {
	// 	panic(err)
	// }
	// ttyHandler.Setup()
	//
	// defer ttyHandler.Close()
	//
	// for data := range ttyHandler.Read() {
	// 	var c [][]tty.Cell
	// 	screen := copy(c, ttyHandler.Screen.Cells)
	// 	_ = screen
	// 	runes := []rune(string(data)) // bytes → runes
	//
	// 	// 建一行 cells
	// 	row := make([]tty.Cell, len(runes))
	// 	for i, r := range runes {
	// 		row[i] = tty.Cell{
	// 			Char:    r,
	// 			Changed: true,
	// 		}
	// 	}
	//
	// 	cells := [][]tty.Cell{row} // 包成二維
	//
	// 	ttyHandler.Screen.ScreenChangeChan <- cells
	// }
	// for {
	// 	// tty.Write("123")
	// 	time.Sleep(1 * time.Second)
	// }
	// kubeConfig := NewKubeConfig()
	//
	// log.Printf("apiversion: %v", kubeConfig.APIVersion)
	// log.Printf("current-context: %v", kubeConfig.CurrentContext)

	// tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	// if err != nil {
	// 	panic(err)
	// }
	// w := bufio.NewWriter(tty)
	// w.WriteString("\033[?1049h")
	// w.Flush()
	//
	// sig := make(chan os.Signal, 1)
	// signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	//
	// w.WriteString("\033[2J\033[H")
	// w.Flush()
	//
	// cleanup := func() {
	// 	w.WriteString("\033[?1049l")
	// 	w.Flush()
	// 	tty.Close()
	// }
	// defer cleanup()
	//
	// go func() {
	// 	<-sig
	// 	cleanup()
	// 	os.Exit(0)
	// }()
	//
	// for {
	// 	w.WriteString("123")
	// 	time.Sleep(1 * time.Second)
	// 	w.Flush()
	// }
}
