package main

import (
	"encoding/gob"
	"log"
	"strings"
	"time"

	"github.com/Distributed-filestorage/p2p"
)

func init() {
	gob.Register(MessageStoreFile{})
	gob.Register(MessageFileContent{})
}

func main() {
	tcpOpts3000 := p2p.TCPTransportOption{
		ListenAddress: "localhost:3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       &p2p.DefaultDecoder{},
	}
	tcpOpts4000 := p2p.TCPTransportOption{
		ListenAddress: "localhost:4000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       &p2p.DefaultDecoder{},
	}
	transportOpts1 := p2p.NewTCPTransport(tcpOpts3000)
	transportOpts2 := p2p.NewTCPTransport(tcpOpts4000)
	fileServerOpts1 := FileServerOpts{
		StorageRoot:       "./data:3000",
		Transport:         transportOpts1,
		PathTransformFunc: CASPathTransformFunc,
		BootstrapNodes:    []string{},
	}
	fileServerOpts2 := FileServerOpts{
		StorageRoot:       "./data:4000",
		Transport:         transportOpts2,
		PathTransformFunc: CASPathTransformFunc,
		BootstrapNodes:    []string{"localhost:3000"},
	}
	fileServer1 := NewFileServer(fileServerOpts1)
	fileServer2 := NewFileServer(fileServerOpts2)
	go func() {
		time.Sleep(2 * time.Second)
		if err := fileServer2.Store(
			"hello.txt",
			strings.NewReader("Hello from node 4000"),
		); err != nil {
			log.Printf("Store failed %v\n", err)
		}
	}()
	transportOpts1.OnPeer = fileServer1.OnPeer
	transportOpts2.OnPeer = fileServer2.OnPeer
	// go func() {
	// 	time.Sleep(5 * time.Second)
	// 	fmt.Println("Stopping server...")
	// 	fileServer.Stop()
	// }()
	go func() {
		if err := fileServer1.Start(); err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(time.Second)
	if err := fileServer2.Start(); err != nil {
		log.Fatal(err)
	}
}
