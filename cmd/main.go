package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	var buf bytes.Buffer

	var w Writer

	w.w = zip.NewWriter(&buf)

	quantity := 10000

	w.wg.Add(quantity)

	for i := 1; i <= quantity; i++ {
		go zipper(i, &w)
	}

	w.wg.Wait()

	err := w.w.Close()
	if err != nil {
		log.Fatal(err)
	}

	bigzip("big", buf.Bytes())
}

type File struct {
	Name string
	Body string
}

func (s *File) NewFile() []byte {
	var buf bytes.Buffer
	zipW := zip.NewWriter(&buf)
	f, err := zipW.Create(fmt.Sprintf("%v.txt", s.Name))
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write([]byte(s.Body))
	if err != nil {
		log.Fatal(err)
	}
	err = zipW.Close()
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

func createfile(name string, body string) []byte {
	f := File{name, body}
	return f.NewFile()
}

func bigzip(name string, buf []byte) {
	err := os.WriteFile(fmt.Sprintf("%v.zip", name), buf, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

type Writer struct {
	w  *zip.Writer
	mu sync.Mutex
	wg sync.WaitGroup
}

func zipper(i int, z *Writer) {
	c := createfile(fmt.Sprintf("%v", i), "content")
	z.NewArchive(c, i)
	z.wg.Done()
}

func (s *Writer) NewArchive(c []byte, i int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := s.w.Create(fmt.Sprintf("%v.zip", i))
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(c)
	if err != nil {
		log.Fatal(err)
	}
}
