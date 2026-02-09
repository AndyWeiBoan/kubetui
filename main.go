package main

import (
	"fmt"
)

func main() {
	test := NewKubeConfig()

	fmt.Println(test.CurrentContext)
}


