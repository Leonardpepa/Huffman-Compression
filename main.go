package main

import (
	"flag"
	"log"
	"os"
)

func main() {

	c := flag.Bool("c", false, "Specify a file to compress")
	d := flag.Bool("d", false, "Specify a file to decompress")

	flag.Parse()

	if *c == false && *d == false {
		log.Fatal(`
	Usage: huffman.exe [OPTIONS]... [FILE] ~1 file
	[OPTIONS]:
	-c FILE ~file to compress
	-d FILE ~file to decompress
`)
	}

	if *c && *d {
		log.Fatal(`
	
	Usage: huffman.exe [OPTIONS]... [FILE] ~1 file
	[OPTIONS]:
	-c FILE ~file to compress
	-d FILE ~file to decompress
`)
	}

	if len(flag.Args()) != 1 {
		log.Fatal(`
	Usage: huffman.exe [OPTIONS]... [FILE]
	[OPTIONS]:
	-c FILE ~file to compress
	-d FILE ~file to decompress
`)
	}

	if *c {
		file := openFile(flag.Args()[0])

		err := Encode(file, "output.hf")

		if err != nil {
			log.Fatal(err)
		}

	} else if *d {
		file := openFile(flag.Args()[0])
		err := Decode(file, "decoded.txt")
		if err != nil {
			log.Fatal(err)
		}
	}

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
