package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestTransformFunc(t *testing.T) {
	key := "myBestPic"
	pathKey := CASPathTransformFunc(key)
	expectedPath := "ff3a4a/df07eb/b0e838/d8bdf4/955913/ea9d81"
	expectedOriginal := "ff3a4adf07ebb0e838d8bdf4955913ea9d81d76f"
	if pathKey.PathName != expectedPath {
		t.Errorf("Expected path: %s, got: %s", expectedPath, pathKey.PathName)
	}
	if pathKey.Original != expectedOriginal {
		t.Errorf("Expected original: %s, got: %s", expectedOriginal, pathKey.Original)
	}
	fmt.Println(pathKey)
}
func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	data := []byte("Some jpeg bytes")
	if err := s.writeStream("myBestPic", bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	r, err := s.Read("myBestPic")
	if err != nil {
		t.Error(err)
	}
	buf, _ := io.ReadAll(r)

	fmt.Println("Read data:", string(buf))
	if string(buf) != string(data) {
		t.Errorf("want %s have %s", data , buf)
	}

}
