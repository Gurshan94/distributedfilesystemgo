package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
    sliceLen := len(hashStr) / blocksize
	path := make([]string, sliceLen)

	for i := 0;i< sliceLen;i++ {
		from, to :=i*blocksize, (i+1)*blocksize
		path[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: strings.Join(path,"/"),
		Filename: hashStr,
	}
	}

type PathKey struct {
	PathName string
	Filename string
}

func (p PathKey) FullPath () string {
	return fmt.Sprintf("%s/%s", p.PathName, p.Filename)
}

func (p PathKey) RootPath () string {
	path := strings.Split((p.PathName),"/")
	if len(path) == 0 {
		return ""
	}
	return path[0]
}

type PathTransformFunc func(string) PathKey

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) delete(key string) error {
	pathKey := s.PathTransformFunc(key)

	defer func () {
		log.Printf("Deleting %s from the disk\n", pathKey.Filename)

	}()

    return os.RemoveAll(pathKey.RootPath())
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.ReadStream(key)
	if err!=nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) ReadStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	return os.Open(pathKey.FullPath())
}

func (s *Store) writeStream(key string, r io.Reader ) error {
	pathKey := s.PathTransformFunc(key)

	if err := os.MkdirAll(pathKey.PathName,os.ModePerm); err!=nil {
		return err
	}

	fullPath := pathKey.FullPath()

	f, err := os.Create(fullPath)
	if err!=nil {
		return err
	}

    defer f.Close()

	n, err := io.Copy(f, r)
	if err!=nil {
		return err
	}

	log.Printf("Wrote %d bytes to %s\n", n, fullPath)
	return nil
}