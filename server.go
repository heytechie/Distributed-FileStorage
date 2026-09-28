package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
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

type Message struct {
	Payload any
}

type MessageStoreFile struct {
	Key  string
	Size int64
}

type MessageFileContent struct {
	Key  string
	Data []byte
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

	return nil
}

func (s *FileServer) broadcast(msg Message) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(msg); err != nil {
		return err
	}
	s.peerLock.Lock()
	peers := make([]p2p.Peer, 0, len(s.peer))
	for _, p := range s.peer {
		peers = append(peers, p)
	}
	s.peerLock.Unlock()
	for _, peer := range peers {
		if err := peer.SendMessage(buf.Bytes()); err != nil {
			return fmt.Errorf("Error sending message to peer %s: %w\n", peer.RemoteAddr(), err)
		}
	}

	return nil
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

func (s *FileServer) Store(key string, r io.Reader) error {
	filebuffer := new(bytes.Buffer)
	tee := io.TeeReader(r, filebuffer)

	if err := s.store.writeStream(key, tee); err != nil {
		return err
	}
	msg := Message{
		Payload: MessageStoreFile{
			Key:  key,
			Size: int64(filebuffer.Len()),
		},
	}

	if err := s.broadcast(msg); err != nil {
		return err
	}

	contentMsg := Message{
		Payload: MessageFileContent{
			Key:  key,
			Data: filebuffer.Bytes(),
		},
	}

	if err := s.broadcast(contentMsg); err != nil {
		return err
	}

	fmt.Printf("Buffered %d bytes for key %s\n", filebuffer.Len(), key)

	return nil
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
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				fmt.Printf("Could not decode the message from %s: %v", rpc.From, err)
				continue
			}

			switch payload := msg.Payload.(type) {
			case MessageStoreFile:
				fmt.Printf("Received store request for key %s of size %d from %s\n", payload.Key, payload.Size, rpc.From)
				// Here you can implement logic to fetch the file from the peer or handle it as needed.
			default:
				fmt.Printf("Received unknown message type from %s\n", rpc.From)
			}
		case <-s.quitch:
			fmt.Println("Server is shutting down")
			return nil

		}
	}
}
