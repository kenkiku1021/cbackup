package main

import (
	"fmt"
	"log"
	"crypto/sha256"
	"io"
	"os"
	"sort"
	"time"
)

type FileDB struct {
	Paths []string
	Size int64
	ModTime time.Time
}

func AppendPath(paths []string, filename string) []string {
	pos := sort.SearchStrings(paths, filename)
	if pos >= len(paths) {
		// do not exist filename
		paths = append(paths, filename)
		sort.Strings(paths)
	}
	return paths
}

func MakeFileHash(filename string) (string, error) {
	h := sha256.New()
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("[warning] Cannot open file (%v)", filename)
		return "", err
	}
	defer f.Close()

	io.Copy(h, f)
	s := fmt.Sprintf("%x", h.Sum(nil))
	return s, nil
}
