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

const defaultRoot = "storeGo"

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
	root              string //folder name of the root directory where the files will be stored
	PathTransformFunc PathTransformFunc
}
type Store struct {
	StoreOpts
}

var DefaultPathTransform = func(key string) PathKey {
	return PathKey{
		PathName: key,
		Original: key,
	}
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransform
	}
	if len(opts.root) == 0 {
		opts.root = defaultRoot
	}
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

func (s *Store) Has(key string) bool {
	path := s.PathTransformFunc(key)
	fullPathWithRoot := filepath.Join(s.root, path.FullPath())
	if _, err := os.Stat(fullPathWithRoot); os.IsNotExist(err) {
		return false
	}
	return true
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
	return os.Open(s.fullPath(key))
}

func (s *Store) fullPath(key string) string {
	path := s.PathTransformFunc(key)
	return filepath.Join(s.root, path.FullPath())
}

func (s *Store) writeStream(key string, r io.Reader) error {
	path := s.PathTransformFunc(key)
	pathAndFilename := filepath.Join(s.root, path.FullPath())

	if err := os.MkdirAll(filepath.Dir(pathAndFilename), os.ModePerm); err != nil {
		return err
	}

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
	fullPath := s.fullPath(key)

	if err := os.Remove(fullPath); err != nil {
		return err
	}

	rootPath := filepath.Join(s.root, path.FirstPathName())
	if err := os.RemoveAll(rootPath); err != nil {
		return err
	}

	fmt.Printf("Deleted file: %s\n", fullPath)
	return nil
}

func (s *Store) Cleanup() error {
	if err := os.RemoveAll(s.root); err != nil {
		return err
	}
	return nil
}
