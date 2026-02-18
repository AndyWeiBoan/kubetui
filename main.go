package main

// "bufio"
// "os"
// "os/signal"
// "syscall"
// "time"

import (
	"fmt"
	"time"

	"kubetui/tty"
)

func main2() {
	fmt.Println(string(1))
}

func main() {
	tty.SetUp()

	defer tty.Close()

	for {
		// tty.Write("123")
		time.Sleep(1 * time.Second)
	}
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
