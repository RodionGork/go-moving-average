package network

import (
	"bufio"
	"log"
	"net"
	"net/textproto"
	"strconv"
	"strings"
	"sync"
)

var listeners sync.WaitGroup

var processor func(string) string

func SetProcessor(p func(string) string) {
	processor = p
}

func StartTCP(params string) {
	port := "1080"
	if params != "" {
		port = params // for now let params contain only port
	}
	if srv, err := net.Listen("tcp", ":"+port); err != nil {
		log.Println("TCP port won't listen, skipping:", err.Error())
	} else {
		listeners.Add(1)
		log.Println("Listening TCP on port", port)
		go serveTCP(srv)
	}
}

func serveTCP(srv net.Listener) {
	defer listeners.Done()
	defer srv.Close()

	for {
		if cli, err := srv.Accept(); err != nil {
			log.Println("TCP accept error:", err.Error())
			continue
		} else {
			go processTCP(cli)
		}
	}
}

func processTCP(client net.Conn) {
	defer client.Close()

	reader := textproto.NewReader(bufio.NewReader(client))
	for {
		line, err := reader.ReadLine()
		if err != nil {
			return
		}
		if res := processor(line); res != "" {
			if _, err := client.Write([]byte(res + "\n")); err != nil {
				log.Println("TCP write error:", err.Error())
			}
		}
	}
}

func StartUDP(params string) {
	port := 1090
	if val, err := strconv.Atoi(params); err == nil {
		port = val
	}
	udpAddr := net.UDPAddr{Port: port, IP: net.ParseIP("0.0.0.0")}
	if udpConn, err := net.ListenUDP("udp", &udpAddr); err != nil {
		log.Println("UDP port won't listen, skipping:", err.Error())
	} else {
		listeners.Add(1)
		log.Println("Listening UDP on port", port)
		go serveUDP(udpConn)
	}
}

func serveUDP(conn *net.UDPConn) {
	defer listeners.Done()
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		if count, addr, err := conn.ReadFrom(buf); err != nil {
			log.Println("UDP read error:", err.Error())
			return
		} else {
			lines := strings.Split(string(buf[0:count]), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				if res := processor(line); res != "" {
					if _, err := conn.WriteTo([]byte(res+"\n"), addr); err != nil {
						log.Println("UDP write error:", err.Error())
					}
				}
			}
		}
	}
}

func WaitAll() {
	listeners.Wait()
}
