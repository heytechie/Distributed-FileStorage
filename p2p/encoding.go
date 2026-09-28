package p2p

import (
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GOBDecoder struct{}

func (d *GOBDecoder) Decode(r io.Reader, msg *RPC) error {
	return gob.NewDecoder(r).Decode(msg)
}

type DefaultDecoder struct{}



const maxMessageSize = 1 << 20 // 1 MB

func (d *DefaultDecoder) Decode(r io.Reader, msg *RPC) error {
	var size uint32
	if err := binary.Read(r, binary.BigEndian, &size); err != nil {
		return err
	}

	if size > maxMessageSize {
		return fmt.Errorf("message too large %d bytes", size)
	}

	msg.Payload = make([]byte, size)

	_, err := io.ReadFull(r, msg.Payload)
	return err

}
