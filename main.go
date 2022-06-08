package main

import (
	"fmt"
	"github.com/wiktor-mazur/dns-go/src/resolver"
	"github.com/wiktor-mazur/dns-go/src/server/udp_server"
	"github.com/wiktor-mazur/dns-go/src/utils"
	"log"
	"net"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	cfg, err := utils.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("error loading config: %s", err.Error()))
	}

	nameResolver := resolver.New(resolver.Config{InternetRootServer: cfg.INTERNET_ROOT_SERVER})

	if cfg.UDP_ENABLED {
		wg.Add(1)
		go func() {
			server := udp_server.New(udp_server.Config{
				ListenIP:   net.ParseIP(cfg.UDP_IP),
				ListenPort: cfg.UDP_PORT,
			}, nameResolver)

			err := server.Start(&wg)
			if err != nil {
				log.Printf("Could not start the UDP server: %s", err.Error())
			}
		}()
	}

	wg.Wait()
}
