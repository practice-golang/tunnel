package main // import "tunnel"

import (
	_ "embed"
	"path/filepath"
	"strings"

	"fmt"
	"log"
	"os"

	"github.com/elliotchance/sshtunnel"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

//go:embed config.yaml
var sampleYAML string

type ProxyInfo struct {
	Address    string `yaml:"address"`
	Port       string `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	PrivateKey string `yaml:"privatekey"`
}

type InternalServerInfo struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
}

type Config struct {
	Proxy          ProxyInfo          `yaml:"proxyserver"`
	InternalServer InternalServerInfo `yaml:"internalserver"`
	LocalPort      string             `yaml:"localport"`
}

func createYAML(iniPath string) {
	if _, err := os.Stat(iniPath); !os.IsNotExist(err) {
		fmt.Printf("File %s already exists.\n", iniPath)
		os.Exit(1)
	}

	f, err := os.Create(iniPath)
	if err != nil {
		log.Fatalln("Create INI: ", err)
	}
	defer f.Close()

	_, err = f.WriteString(sampleYAML)
	if err != nil {
		log.Fatalln("Create INI: ", err)
	}

	fmt.Println(iniPath + " is created")
	fmt.Println("Please modify " + iniPath + " then run again")

	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("* Run:        tunnel ./config.yaml")
		fmt.Println("* Get config: tunnel -getyaml")
		os.Exit(1)
	}

	if os.Args[1] == "-getyaml" {
		createYAML("config_sample.yaml")
	}

	configFile := os.Args[1]

	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Println("Error when " + configFile + " reading")
		os.Exit(1)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error when " + configFile + " parsing")
		os.Exit(1)
	}

	proxyAddr := config.Proxy.Username + "@" + config.Proxy.Address + ":" + config.Proxy.Port
	proxyPw := config.Proxy.Password
	inAddr := config.InternalServer.Address + ":" + config.InternalServer.Port

	authWay := ssh.Password(proxyPw)
	if config.Proxy.PrivateKey != "" {
		pemPath := config.Proxy.PrivateKey
		if !filepath.IsAbs(pemPath) {
			if strings.HasPrefix(pemPath, "~/") {
				dirname, _ := os.UserHomeDir()
				pemPath = filepath.Join(dirname, pemPath[2:])
			}

			pemPath, _ = filepath.Abs(pemPath)
			fmt.Println("private key path:", pemPath)
		}

		authWay = sshtunnel.PrivateKeyFile(pemPath)
		fmt.Println("Use private key instead password")
	}

	tunnel, err := sshtunnel.NewSSHTunnel(proxyAddr, authWay, inAddr, config.LocalPort)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	tunnel.Start()
}
