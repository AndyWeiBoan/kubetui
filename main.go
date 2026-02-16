package main

import (
	_"fmt"
	"log"
)

func main() {
	kubeConfig := NewKubeConfig()

	log.Printf("apiversion: %v", kubeConfig.APIVersion)
	log.Printf("current-context: %v", kubeConfig.CurrentContext)

}


