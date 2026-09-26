package main

import (
	"fmt"
	"sync"

	"github.com/Distributed-filestorage/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts
	store    *Store
	quitch   chan struct{}
	peer     map[string]p2p.Peer
	peerLock sync.Mutex
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
		peer:           make(map[string]p2p.Peer),
	}
}
func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	addr := p.RemoteAddr().String()
	s.peerLock.Lock()
	// defer s.peerLock.Unlock()
	s.peer[addr] = p
	s.peerLock.Unlock()
	fmt.Printf("New Peer connected: %s\n", addr)
	return p.Send([]byte("Helloooo"))
}

func (s *FileServer) bootstrapNetwork() {
	for _, addr := range s.BootstrapNodes {
		if addr == "" {
			continue
		}

		go func(addr string) {
			fmt.Printf("Atthemting to connect to %s\n", addr)
			if err := s.Transport.Dial(addr); err != nil {
				fmt.Printf("Could not connect to %s : %v\n", addr, err)
			}
		}(addr)
	}
}

func (s *FileServer) Start() error {
	// Start the transport to listen for incoming connections
	err := s.Transport.ListenAndAccept()
	if err != nil {
		return err
	}
	defer s.Transport.Close()
	fmt.Println("Server is listening")
	s.bootstrapNetwork()
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
