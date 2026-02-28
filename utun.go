package utun

import (
	"io"
	"log"
	"net"
	"net/netip"
	"sync/atomic"
	"unsafe"
)

func Server(tun io.ReadWriter, c *net.UDPConn, key []byte) {
	var cAddr atomic.Pointer[netip.AddrPort]

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
				xor2(b, key)
				_, err := c.WriteToUDPAddrPort(b, *a)
				if err != nil {
					log.Println("WriteTo err:", err)
					cAddr.Store(nil)
				}
			}
		}
	}()

	buf := make([]byte, 1500)
	for {
		n, addr, err := c.ReadFromUDPAddrPort(buf)
		if err != nil {
			log.Println("ReadFrom err:", err)
		}

		if n == 0 {
			continue
		}

		b := buf[:n]

		xor2(b, key)

		if old := cAddr.Load(); old == nil || *old != addr {
			cAddr.Store(&addr)
		}

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

			xor2(b, key)

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
		xor2(b, key)

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
	for i := range data {
		data[i] ^= key[j]
		j += 1
		if j >= len(key) {
			j = 0
		}
	}
}

// 注意：这里假设 key 长度是 8 的倍数。
func xor2(data, key []byte) {
	dataLen := len(data)
	if dataLen == 0 {
		return
	}

	keyLen := len(key)
	if keyLen == 0 {
		return
	}
	mask := keyLen - 1

	i := 0

	for ; i <= dataLen-8; i += 8 {
		k64 := *(*uint64)(unsafe.Pointer(&key[i&mask]))
		*(*uint64)(unsafe.Pointer(&data[i])) ^= k64
	}

	for ; i < dataLen; i++ {
		data[i] ^= key[i&mask]
	}
}
