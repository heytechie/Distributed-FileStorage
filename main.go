package main

import (
	"fmt"
	"log"

	"github.com/Distributed-filestorage/p2p"
)

func main() {
	tr := p2p.NewTCPTransport(":3000")

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatalf(fmt.Sprintf("Error starting TCP transport: %s\n", err))
	}
	select {}
}
