package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestDeleteFunc(t *testing.T) {
	s := newStore()
	defer func() {
		// time.Sleep(2*time.Second)
		tearDownStore(s, t)
	}()
	key := "heheheheh"
	data := []byte("Some jpeg bytes")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	s.root = "storeGo" // Set the root directory
	// time.Sleep(3 * time.Second) // Ensure the file is written before deletion
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

	s := newStore()
	defer tearDownStore(s, t) // Ensure cleanup after the test

	for i := 0; i < 50000; i++ {
		fmt.Println("file", i)
		data := []byte(fmt.Sprintf("Some jpeg bytes %d", i))
		key := fmt.Sprintf("fileassa%d", i)
		if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}
		if has := s.Has(key); !has {
			t.Errorf("Expected file to exist, but it doesn't")
		}
		r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}
		buf, _ := io.ReadAll(r)

		fmt.Println("Read data:", string(buf))
		if string(buf) != string(data) {
			t.Errorf("want %s have %s", data, buf)
		}
		s.Delete(key)
	}

}

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(opts)
}

func tearDownStore(s *Store, t *testing.T) {
	// Clean up the store directory after tests
	if err := s.Cleanup(); err != nil {
		t.Error("Error cleaning up store directory:", err)
	}
}
