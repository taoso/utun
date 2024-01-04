package utun

import (
	"io"
	"log"
	"net"
	"sync/atomic"
)

type PacketConn interface {
	ReadFrom(p []byte) (n int, addr net.Addr, err error)
	WriteTo(p []byte, addr net.Addr) (n int, err error)
}

func Server(tun io.ReadWriter, c PacketConn, key []byte) {
	var cAddr atomic.Value

	go func() {
		buf := make([]byte, 1500)
		for {
			n, err := tun.Read(buf)
			if err != nil {
				log.Println("tun read err:", err)
			}

			if n == 0 {
				continue
			}

			b := buf[:n]

			if a := cAddr.Load(); a != nil {
				xor(b, key)
				_, err := c.WriteTo(b, a.(net.Addr))
				if err != nil {
					log.Println("WriteTo err:", err)
					cAddr.Store(nil)
				}
			}
		}
	}()

	for {
		buf := make([]byte, 1500)
		n, addr, err := c.ReadFrom(buf)
		if err != nil {
			log.Println("ReadFrom err:", err)
		}

		if n == 0 {
			continue
		}

		b := buf[:n]

		xor(b, key)

		cAddr.Store(addr)

		if _, err := tun.Write(b); err != nil {
			log.Println("tun write err:", err)
		}
	}
}

func Client(tun, conn io.ReadWriter, key []byte) {
	go func() {
		buf := make([]byte, 1500)
		for {
			n, err := tun.Read(buf)

			if err != nil {
				log.Println("tun read err:", err)
			}

			if n == 0 {
				continue
			}

			b := buf[:n]

			xor(b, key)

			if _, err := conn.Write(b); err != nil {
				log.Println("UDP write err:", err)
			}
		}
	}()

	buf := make([]byte, 1500)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("UDP read err:", err)
		}

		if n == 0 {
			continue
		}

		b := buf[:n]
		xor(b, key)

		if _, err := tun.Write(b); err != nil {
			log.Println("tun write err:", err)
		}
	}
}

func xor(data, key []byte) {
	if len(key) == 0 {
		return
	}
	j := 0
	for i := 0; i < len(data); i++ {
		data[i] ^= key[j]
		j += 1
		if j >= len(key) {
			j = 0
		}
	}
}
