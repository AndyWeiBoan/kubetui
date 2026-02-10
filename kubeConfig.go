package main

import (
	"bufio"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"strings"
)

type KubeConfig struct{
	ApiVersion string
	Kind string
	Contexts []Context
	CurrentContext string
}

type Option func(*KubeConfig)

type Context struct {
	Cluster Clusters
	User Users
}

type Clusters struct {
	CertificateAuthorityData string
	Server string
	Name string
}

type Users struct {
	User User
	Name string
}

type User struct {
	ClientCertificateData string
	ClientKeyData string
}

func NewKubeConfig() *KubeConfig{

	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".kube", "config")
	for line := range readKubeConfigFile(path) {
		item := strings.Split(line, ":")
		fmt.Println(item[0])
	}
	return &KubeConfig {

	}
}

func readKubeConfigFile(path string) iter.Seq[string] {
	return func(yield func(string) bool) {
		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if !yield(line) {
				return
			}
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}
}

