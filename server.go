package main

import (
	"fmt"

	"github.com/Distributed-filestorage/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
}

type FileServer struct {
	FileServerOpts
	store  *Store
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	store := NewStore(storeOpts)
	return &FileServer{
		FileServerOpts: opts,
		store:          store,
		quitch:         make(chan struct{}),
	}
}
func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) Start() error {
	// Start the transport to listen for incoming connections
	err := s.Transport.ListenAndAccept()
	if err != nil {
		return err
	}
	defer s.Transport.Close()
	fmt.Println("Server is listening")
	for {
		select {
		case rpc := <-s.Transport.Consume():
			fmt.Printf("From %s: %s\n", rpc.From, rpc.Payload)
		case <-s.quitch:
			fmt.Println("Server is shutting down")
			return nil

		}
	}
}
