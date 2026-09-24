package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestDeleteFunc(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "myBestPic"
	data := []byte("Some jpeg bytes")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	s.root = "storeGo"          // Set the root directory
	time.Sleep(3 * time.Second) // Ensure the file is written before deletion
	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
	// Try to read the deleted file
	_, err := s.Read(key)
	if err == nil {
		t.Errorf("Expected error when reading deleted file, got nil")
	} else {
		fmt.Println("Successfully deleted the file and confirmed it cannot be read.")
	}
}

func TestTransformFunc(t *testing.T) {
	key := "myBestPic"
	pathKey := CASPathTransformFunc(key)
	expectedPath := "ff3a4/adf07/ebb0e/838d8/bdf49/55913/ea9d8/1d76f"
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
	data := []byte("Some jpeg bytes is one")
	if err := s.writeStream("myBestPic", bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	// r, err := s.Read("myBestPic")
	// if err != nil {
	// 	t.Error(err)
	// }
	// buf, _ := io.ReadAll(r)

	// fmt.Println("Read data:", string(buf))
	// if string(buf) != string(data) {
	// 	t.Errorf("want %s have %s", data, buf)
	// }
}
