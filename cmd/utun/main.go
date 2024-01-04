package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	_ "net/http/pprof"

	"github.com/songgao/water"
	"github.com/taoso/utun"
)

func main() {
	var listen string
	var connect string
	var skey string
	var pprof string

	flag.StringVar(&listen, "listen", "", "listen address")
	flag.StringVar(&connect, "connect", "", "connect address")
	flag.StringVar(&skey, "key", "", "shared xor key")
	flag.StringVar(&pprof, "pprof", "", "pprof address")

	flag.Parse()

	if skey == "" {
		log.Fatal("shared xor key can not be empty")
	}
	key, err := hex.DecodeString(skey)
	if err != nil {
		log.Fatal(err)
	}

	cfg := water.Config{DeviceType: water.TUN}
	tun, err := water.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer tun.Close()

	fmt.Println("tun:", tun.Name())

	if pprof != "" {
		go http.ListenAndServe(pprof, nil)
	}

	if listen != "" {
		c, err := net.ListenPacket("udp", listen)
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()

		utun.Server(tun, c, key)
	} else if connect != "" {
		c, err := net.Dial("udp", connect)
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()

		utun.Client(tun, c, key)
	} else {
		log.Fatal("you must set either listen or connect")
	}
}
