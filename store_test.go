package main

import (
	"bytes"
	"io/ioutil"
	"testing"
	"fmt"
)
func TestCASPathTransformFunc(t *testing.T) {
	key := "picture"
	pathKey := CASPathTransformFunc(key)
	expectedOriginalKey := "0e66b49d04a43c19bcb2e13a53aff5d174f610fc"
	expectedFilepath := "0e66b/49d04/a43c1/9bcb2/e13a5/3aff5/d174f/610fc"
	if pathKey.PathName != expectedFilepath {
		t.Errorf("Expected %s, got %s", expectedFilepath, pathKey.PathName)
	}
	if pathKey.Filename != expectedOriginalKey {
		t.Errorf("Expected %s, got %s", expectedFilepath, pathKey.PathName)
	}
}




func TestStoreDeleteKey(t *testing.T) {
	StoreOpts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s:=NewStore(StoreOpts)
	key:="secret"
	data := []byte("secret data")

	if err:=s.writeStream(key,bytes.NewReader(data)); err!=nil {
		t.Error(err)
	}

	if err := s.delete(key);err!=nil {
		t.Error(err)
	}
}
func TestStore(t *testing.T) {
	StoreOpts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	s:=NewStore(StoreOpts)
	key:="hash"
	data := []byte("Hello World")

	if err:=s.writeStream(key,bytes.NewReader(data)); err!=nil {
		t.Error(err)
	}

	r, err := s.Read(key)
	if err!=nil {
		t.Error(err)
	}
	b, _ := ioutil.ReadAll(r)

    fmt.Println(string(b))

	if string(b) != string(data) {
		t.Errorf("Expected %s, got %s", string(data), string(b))
	}

	s.delete(key)
}