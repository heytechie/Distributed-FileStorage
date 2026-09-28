package main

import (
	"bytes"
	"encoding/gob"
	"testing"
)

func TestMessageGobRoundTrip(t *testing.T) {
	original := Message{
		Payload: MessageStoreFile{
			Key:  "test.txt",
			Size: 20,
		},
	}

	var wire bytes.Buffer
	if err := gob.NewEncoder(&wire).Encode(original); err != nil {
		t.Fatalf("Failed to encode: %v", err)
	}

	var decoded Message
	if err := gob.NewDecoder(&wire).Decode(&decoded); err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}

	got, ok := decoded.Payload.(MessageStoreFile)
	if !ok {
		t.Fatalf("Decoded payload is not of type MessageStoreFile")
	}

	if got.Key != "test.txt" || got.Size != 20 {
		t.Errorf("Decoded payload does not match original. Got: %+v, Want: %+v", got, original.Payload)
	}

}
