package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type KubeConfig struct{
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

	kubeConfig, err := loadKubeConfig()

	if (err != nil) {
		panic(err)
	}

	return &KubeConfig {
		CurrentContext:  kubeConfig,
	}
}

func loadKubeConfig() (string,error) {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".kube", "config")
	bytes, err := os.ReadFile(path)
	if err != nil {

		fmt.Println(err)
		return "", err
	}
	result := string(bytes)
	return result, nil
}

