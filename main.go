package main

import (
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const overrideCount = 7

func main() {
	paths := os.Args[1:]
	for _, path := range paths {
		destroy(path)
	}
}

func destroy(path string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		os.Stderr.WriteString("cannot open file '" + path + "': " + err.Error())
		return
	}

	size, err := f.Seek(0, io.SeekEnd)
	check(err)
	for i := 0; i < overrideCount; i++ {
		_, err = f.Seek(0, io.SeekStart)
		check(err)
		_, err = io.CopyN(f, r, size)
		check(err)
		check(f.Sync())
	}
	check(f.Close())
	f, err = os.OpenFile(path, os.O_APPEND, 0)
	check(err)
	_, err = f.Seek(0, io.SeekEnd)
	check(err)
	_, err = io.CopyN(f, r, 100+int64(r.Intn(100)))
	check(err)
	check(f.Close())

	equalLenRandPath := filepath.Join(
		filepath.Dir(path),
		randFilename(r, len(filepath.Base(path))),
	)
	check(os.Rename(path, equalLenRandPath))

	randPath := filepath.Join(
		filepath.Dir(path),
		randFilename(r, 5+r.Intn(20)),
	)
	check(os.Rename(equalLenRandPath, randPath))

	check(os.Remove(randPath))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func randFilename(r *rand.Rand, count int) string {
	s := make([]byte, count)
	for i := range s {
		s[i] = fileChars[rand.Intn(len(fileChars))]
	}
	return string(s)
}

const fileChars = "abcdefghijklmnopqrstuvwxyz0123456789"
