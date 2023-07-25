package main

import (
	"os"

	"github.com/rodiongork/go-moving-average/pkg/counter"
	"github.com/rodiongork/go-moving-average/pkg/network"
)

func main() {
	network.SetProcessor(counter.RequestProcessor)
	network.StartTCP(os.Getenv("TCP_PORT"))
	network.StartUDP(os.Getenv("UDP_PORT"))
	network.WaitAll()
}
