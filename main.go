package main

import (
	"flag"
	"huffmanCompression/huffman"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	c := flag.Bool("c", false, "Specify a file to compress")
	d := flag.Bool("d", false, "Specify a file to decompress")
	h := flag.Bool("h", false, "Usage")

	flag.Parse()

	if *h {
		usage()
	}

	if *c == false && *d == false {
		usage()
	}

	if *c && *d {
		log.Fatal("You can choose either compress or decompress. Try '--help' for more information.")
	}

	if len(flag.Args()) != 1 {
		log.Fatal("Please provide a file. Try '--help' for more information.")
	}

	if *c {
		file := openFile(flag.Args()[0])

		err := huffman.Encode(file, filepath.Base(file.Name())+".hf")

		if err != nil {
			log.Fatal(err)
		}

	} else if *d {
		file := openFile(flag.Args()[0])
		err := huffman.Decode(file, strings.Split(filepath.Base(file.Name()), ".hf")[0])
		if err != nil {
			log.Fatal(err)
		}
	}

}

func usage() {
	log.Fatal(`
Usage: huffman.exe [OPTIONS]... [FILE]
[OPTIONS]:
  -c FILE ~file to compress
  -d FILE ~file to decompress
[FILE]:
  the file specified is the file to compress or decompress,
  when compressed the file will be saved with a .hf extension,
  when decompressed the file will be saved as a copy of the original file
`)
}

// TODO create cli api
func openFile(filename string) *os.File {
	file, err := os.Open(filename)
	log.Println("Opening file... ", filename)
	if err != nil {
		log.Fatal(err)
	}

	return file
}
