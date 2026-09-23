package p2p

import (
	"fmt"
	"net"
)

type TCPPeer struct {
	conn net.Conn
	//basically if we dial and retrive then its a outbound connection, if we accept then its inbound connection
	outbound bool
}

func NewPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOption struct {
	ListenAddress string
	HandshakeFunc HandShakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOption
	listner net.Listener
	rpcChan chan RPC
}

func NewTCPTransport(opts TCPTransportOption) *TCPTransport {
	return &TCPTransport{
		TCPTransportOption: opts,
		rpcChan:            make(chan RPC),
	}
}

func (p *TCPPeer) Close() error {
	err := p.conn.Close()
	return err
}

// Consume implements the Transport interface and returns a channel of RPC messages recieved from peers. It allows the application to consume incoming RPC messages.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcChan

}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listner, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()

	return nil

}

func (t *TCPTransport) startAcceptLoop() {
	for {

		conn, err := t.listner.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
			continue
		}
		fmt.Printf("Incomming new connection:%+v\n", conn)
		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewPeer(conn, true)
	var err error
	defer func() {
		fmt.Printf("Dropping Peer connection %s\n", err)
		conn.Close()
	}()
	if err := t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		err = t.OnPeer(peer)
		if err != nil {
			return
		}
	}
	rpc := RPC{}
	for {
		err = t.Decoder.Decode(conn, &rpc)
		if err != nil {
			fmt.Printf("Error decoding RPC: %s\n", err)
			return
		}
		rpc.From = conn.RemoteAddr()
		// fmt.Printf("Recieved from %s\n", rpc.From)
		// fmt.Printf("Message: %s\n", string(rpc.Payload))
		t.rpcChan <- rpc
	}
}
