package main

import (
	"bufio"
	"fmt"
	"iter"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type KubeConfig struct {
	APIVersion     string
	Kind           string
	Contexts       []Context
	CurrentContext string
}

type Option func(*KubeConfig)

type ContextRaw struct {
	Cluster string
	User    string
}

type Context struct {
	Cluster Cluster
	User    User
}

type Cluster struct {
	CertificateAuthorityData string

	Server string
	Name   string
}

type User struct {
	ClientKeyData         string
	ClientCertificateData string
	Token                 string
	Name                  string
}

func NewKubeConfig() *KubeConfig {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".kube", "config")
	kubeConfig := &KubeConfig{
		Contexts:   make([]Context, 0),
		APIVersion: "",
		Kind:       "",
	}
	content := loadKubeConfig(path)
	lines := strings.Split(content, "\n")

	apiversion, err := ParseAPIVersion(lines)
	if err != nil {
		panic(err)
	}
	kubeConfig.APIVersion = apiversion

	currentContext, err := ParseCurrentContext(lines)
	if err != nil {
		panic(err)
	}
	kubeConfig.CurrentContext = currentContext

	clusters, err := ParsseClusters(lines)
	if err != nil {
		panic(err)
	}

	users, err := ParseUsers(lines)
	if err != nil {
		panic(err)
	}

	contextRaws, err := ParseContextRaws(lines)
	if err != nil {
		panic(err)
	}

	for _, contextRaw := range contextRaws {
		cluster, clusterExists := clusters[contextRaw.Cluster]
		user, userExists := users[contextRaw.User]
		if clusterExists && userExists {
			kubeConfig.Contexts = append(kubeConfig.Contexts, Context{
				Cluster: *cluster,
				User:    *user,
			})
		} else {
			if !clusterExists {
				log.Printf("the cluster not exitsts: %v", contextRaw.Cluster)
			}

			if !userExists {
				log.Printf("the user not exitsts: %v", contextRaw.User)
			}
		}
	}

	return kubeConfig
}

func ParseAPIVersion(lines []string) (string, error) {
	for _, line := range lines {
		if strings.Contains(line, "apiVersion:") {
			_, value, found := strings.Cut(line, ":")
			if found {
				return value, nil
			}
		}
	}
	return "", fmt.Errorf("ApiVersion not found while kubeconfig tokenization")
}

func ParseCurrentContext(lines []string) (string, error) {
	for _, line := range lines {
		if strings.HasPrefix(line, "current-context:") {
			_, value, found := strings.Cut(line, ":")

			if found {
				return value, nil
			}
			return "", fmt.Errorf("the current-context invalid format: %v ", line)
		}
	}
	return "", fmt.Errorf("the current-context not found")
}

func ParsseClusters(lines []string) (map[string]*Cluster, error) {
	result := make(map[string]*Cluster)
	var current *Cluster
	for index, line := range lines {
		if strings.Contains(line, "clusters:") {
			if index == len(lines)-1 {
				log.Printf("invalid kubeconfig format, clusters should not be last line")
				break
			}
			if strings.Contains(line, "[") || strings.Contains(lines[index+1], "[") {
				// flow-style
			} else {
				// block-style
				for i := index + 1; i < len(lines); i++ {
					key, value, found := strings.Cut(lines[i], ":")
					hasIndent := strings.HasPrefix(key, " ")
					hasArrayToken := strings.Contains(key, "-")
					if !hasIndent && !hasArrayToken {
						break
					}
					if hasArrayToken {
						current = &Cluster{}
					}
					if found {
						if strings.Contains(key, "name") {
							current.Name = value
							result[current.Name] = current
						} else if strings.Contains(key, "server") {
							current.Server = value
						} else if strings.Contains(key, "certificate-authority-data") {
							current.CertificateAuthorityData = value
						}
					}
				}
			}
			break
		}
	}
	return result, nil
}

func ParseUsers(lines []string) (map[string]*User, error) {
	result := make(map[string]*User)
	var current *User
	for index, line := range lines {
		if strings.Contains(line, "users:") {
			if index == len(lines)-1 {
				log.Printf("invalid kubeconfig format, users should not be last line")
				break
			}
			if strings.Contains(line, "[") || strings.Contains(lines[index+1], "[") {
				// flow-style
			} else {
				// block-style
				for i := index + 1; i < len(lines); i++ {
					key, value, found := strings.Cut(lines[i], ":")
					hasIndent := strings.HasPrefix(key, " ")
					hasArrayToken := strings.Contains(key, "-")
					if !hasIndent && !hasArrayToken {
						break
					}
					if hasArrayToken {
						current = &User{}
					}
					if found {
						if strings.Contains(key, "name") {
							current.Name = value
							result[current.Name] = current
						} else if strings.Contains(key, "client-certificate-data") {
							current.ClientCertificateData = value
						} else if strings.Contains(key, "client-key-data") {
							current.ClientKeyData = value
						} else if strings.Contains(key, "toke") {
							current.Token = value
						}
					}
				}
			}
			break
		}
	}
	return result, nil
}

func ParseContextRaws(lines []string) ([]*ContextRaw, error) {
	result := make([]*ContextRaw, 0)
	var current *ContextRaw
	for index, line := range lines {
		if strings.Contains(line, "contexts:") {
			if index == len(lines)-1 {
				log.Printf("invalid kubeconfig format, users should not be last line")
				break
			}
			if strings.Contains(line, "[") || strings.Contains(lines[index+1], "[") {
				// flow-style
			} else {
				// block-style
				for i := index + 1; i < len(lines); i++ {
					key, value, found := strings.Cut(lines[i], ":")
					hasIndent := strings.HasPrefix(key, " ")
					hasArrayToken := strings.Contains(key, "-")
					if !hasIndent && !hasArrayToken {
						break
					}
					if hasArrayToken {
						current = &ContextRaw{}
					}
					if found {
						if strings.Contains(key, "cluster") {
							current.Cluster = value
							result = append(result, current)
						} else if strings.Contains(key, "user") {
							current.User = value
						}
					}
				}
			}
			break
		}
	}
	return result, nil
}

func (kubeconfig *KubeConfig) SetAPIVersion(text string) *KubeConfig {
	_, value, found := strings.Cut(text, ":")

	if !found {
		log.Panicf("Invalid KubeConfig format: %s", text)
	}
	kubeconfig.APIVersion = strings.Trim(value, "")
	log.Printf("set apiVerion: %s", kubeconfig.APIVersion)
	return kubeconfig
}

func loadKubeConfig(path string) string {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(bytes)
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
