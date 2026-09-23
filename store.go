package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i+1)*blockSize
		paths[i] = hashStr[from:to]
	}
	return PathKey{
		PathName: strings.Join(paths, "/"),
		Original: hashStr,
	}
}

type PathKey struct {
	PathName string
	Original string
}

func (p PathKey) FullPath() string {
	return filepath.Join(p.PathName, p.Original)
}

type PathTransformFunc func(string) PathKey
type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}
type Store struct {
	StoreOpts
}

var DefaultPathTransform = func(key string) string {
	return key
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		panic("PathName is empty")
	}
	return paths[0]
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	path := s.PathTransformFunc(key)
	return os.Open(path.FullPath())
}

func (s *Store) writeStream(key string, r io.Reader) error {
	path := s.PathTransformFunc(key)

	if err := os.MkdirAll(path.PathName, os.ModePerm); err != nil {
		return err
	}

	pathAndFilename := path.FullPath()

	f, err := os.Create(pathAndFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	log.Printf("Wrote %d bytes to %s", n, pathAndFilename)

	return nil
}

func (s *Store) Delete(key string) error {
	path := s.PathTransformFunc(key)
	fullPath := path.FullPath()

	if err := os.RemoveAll(path.FirstPathName()); err != nil {
		return err
	}

	fmt.Printf("Deleted file: %s\n", fullPath)
	return nil
}
