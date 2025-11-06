package main

import (
	"compress/gzip"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: dumphex <file>")
		os.Exit(1)
	}
	p := os.Args[1]
	f, err := os.Open(p)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b := make([]byte, 2)
	_, err = f.Read(b)
	if err != nil {
		panic(err)
	}
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}
	var data []byte
	// gzip magic is 0x1f 0x8b (big-endian when combined)
	if (int(b[0])<<8)|int(b[1]) == 0x1f8b {
		gr, err := gzip.NewReader(f)
		if err != nil {
			panic(err)
		}
		defer gr.Close()
		data, err = io.ReadAll(gr)
		if err != nil {
			panic(err)
		}
	} else {
		data, err = io.ReadAll(f)
		if err != nil {
			panic(err)
		}
	}

	n := 64
	if len(data) < n {
		n = len(data)
	}
	fmt.Println(hex.Dump(data[:n]))
}
