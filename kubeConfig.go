package main

import (
	"bufio"
	"iter"
	"os"
	"log"
	"path/filepath"
	"strings"
)

type KubeConfig struct{
	APIVersion string
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
	kubeConfig := &KubeConfig{}
	for line := range readKubeConfigFile(path) {
		if strings.Contains(line, "apiVersion:") {
			kubeConfig.SetAPIVersion(line)
		}
	}
	return kubeConfig
}

func (kubeconfig * KubeConfig) SetAPIVersion(text string) *KubeConfig {

	_, value, found := strings.Cut(text, ":")

	if (!found) {
		log.Panicf("Invalid KubeConfig format: %s", text)
	}
	kubeconfig.APIVersion = strings.Trim(value, "")
	log.Printf("set apiVerion: %s", kubeconfig.APIVersion)
	return kubeconfig
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

