package main

import (
	"fmt"
	"io"
	"log"

	"github.com/k1LoW/ghfs"
)

func main() {
	fsys, err := ghfs.New("golang", "time")
	if err != nil {
		log.Fatal(err)
	}
	f, err := fsys.Open("README.md")
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)
}
