package main

import (
	"fmt"
	"log"

	"github.com/Distributed-filestorage/p2p"
)

func Onpeer(peer p2p.Peer) error {
	peer.Close()
	fmt.Printf("New Peer connected: %s\n", peer)
	return nil
}
func main() {
	tcpOpts := p2p.TCPTransportOption{
		ListenAddress: "localhost:3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       &p2p.DefaultDecoder{},
		OnPeer:        Onpeer,
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("Recieved from %s\n", msg.From)
			fmt.Printf("Message: %s\n", string(msg.Payload))
		}
	}()
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}
}
