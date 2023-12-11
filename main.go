package main // import "tunnel"

import (
	"log"
	"os"
	"time"

	"github.com/elliotchance/sshtunnel"
	"golang.org/x/crypto/ssh"
)

type ProxyInfo struct {
	Address  string
	Port     string
	Username string
	Password string
}

type InternalServerInfo struct {
	Address string
	Port    string
}

type Config struct {
	Proxy          ProxyInfo
	InternalServer InternalServerInfo
	OutsidePort    string
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: tunnel config.yaml")
	}

	configFile := os.Args[1]

	log.Println(configFile)
	return

	config := Config{
		ProxyInfo{"192.168.0.1", "22", "proxy_user_id", "proxy_user_password"},
		InternalServerInfo{"192.168.1.2", "22"},
		"1222",
	}

	proxyAddr := config.Proxy.Username + "@" + config.Proxy.Address + ":" + config.Proxy.Port
	proxyPw := config.Proxy.Password
	inAddr := config.InternalServer.Address + ":" + config.InternalServer.Port

	tunnel, err := sshtunnel.NewSSHTunnel(proxyAddr, ssh.Password(proxyPw), inAddr, config.OutsidePort)
	if err != nil {
		log.Fatal(err)
	}

	tunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)

	go tunnel.Start()
	for {
		time.Sleep(100 * time.Second)
	}
}
