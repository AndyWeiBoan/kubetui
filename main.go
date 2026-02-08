package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	print_kubeConfig()
}

func print_kubeConfig() {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".kube", "config")

	bytes, err := os.ReadFile(path)
	if err != nil {

		fmt.Println(err)
		return
	}
	println(string(bytes))
}
