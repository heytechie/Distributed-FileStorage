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
		// OnPeer:        Onpeer,
	}
	transportOpts := p2p.NewTCPTransport(tcpOpts)
	fileServerOpts := FileServerOpts{
		StorageRoot:       "./data:3000",
		Transport:         transportOpts,
		PathTransformFunc: CASPathTransformFunc,
	}
	fileServer := NewFileServer(fileServerOpts)
	// go func() {
	// 	time.Sleep(5 * time.Second)
	// 	fmt.Println("Stopping server...")
	// 	fileServer.Stop()
	// }()
	if err := fileServer.Start(); err != nil {
		log.Fatal(err)
	}
}
