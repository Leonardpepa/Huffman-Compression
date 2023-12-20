package main

import (
	"log"
	"os"
)

func main() {
	file := readFile("input/gutenberg.txt")
	err := Encode(file, "output.hf")
	if err != nil {
		log.Fatal(err)
	}

	file = readFile("output.hf")
	err = Decode(file)

	if err != nil {
		log.Fatal(err)
	}
}

// TODO create cli api
func readFile(filename string) *os.File {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	return file
}
