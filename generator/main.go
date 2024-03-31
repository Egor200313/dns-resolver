package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
)

const charset = "abcdefghijklmnopqrstuvwxyz"

func randomARecord() string {
	l := 10
	b := make([]byte, l)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return fmt.Sprintf("%s.com.\tIN\tA\t%d.%d.%d.%d", b, rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

func main() {
	f, err := os.Create("random_records.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for range 10000 {
		line := randomARecord()
		_, err := io.WriteString(f, line+"\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	line := randomARecord()
	_, err = io.WriteString(f, line)
	if err != nil {
		log.Fatal(err)
	}
}
