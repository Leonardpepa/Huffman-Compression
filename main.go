package main

import (
	"flag"
	"huffmanCompression/huffman"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var programName = filepath.Base(os.Args[0])

func main() {

	compress := flag.Bool("c", false, "Specify a file to compress")
	decompress := flag.Bool("d", false, "Specify a file to decompress")
	help := flag.Bool("h", false, "Usage")

	flag.Parse()

	if *help {
		usage()
	}

	if *compress == false && *decompress == false {
		usage()
	}

	if *compress && *decompress {
		log.Fatal("You can choose either compress or decompress. Try '--help' for more information.")
	}

	if len(flag.Args()) != 1 {
		log.Fatal("Please provide a file. Try '--help' for more information.")
	}

	if *compress {
		file := openFile(flag.Args()[0])

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Println(err)
			}
		}(file)

		// add the .hf extension
		outputPath := filepath.Base(file.Name()) + ".hf"
		err := huffman.Encode(file, outputPath)

		if err != nil {
			log.Fatal(err)
		}

	} else if *decompress {
		file := openFile(flag.Args()[0])

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Println(err)
			}
		}(file)

		// remove the .hf extension
		outputPath := strings.Split(filepath.Base(file.Name()), ".hf")[0]
		err := huffman.Decode(file, outputPath)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func usage() {
	log.Fatalf(`
Usage: %s [OPTIONS]... [FILE]
[OPTIONS]:
  -c FILE ~file to compress
  -d FILE ~file to decompress
[FILE]:
  the file specified is the file to compress or decompress,
  when compressed the file will be saved with a .hf extension,
  when decompressed the file will be saved as a copy of the original file
`, programName)
}

func openFile(filename string) *os.File {
	file, err := os.Open(filename)
	log.Println("Opening file... ", filename)
	if err != nil {
		log.Fatal(err)
	}

	return file
}
