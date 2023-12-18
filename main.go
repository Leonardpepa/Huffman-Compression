package main

import (
	"log"
	"os"
)

func main() {
	//charFrequencies := map[rune]uint64{
	//	'Z': 2,
	//	'K': 7,
	//	'M': 24,
	//	'C': 32,
	//	'U': 37,
	//	'D': 42,
	//	'L': 42,
	//	'E': 120,
	//}
	//
	//root := getHuffmanTreeFromFrequencies(charFrequencies)

	root := getHuffmanTreeFromFile("input/gutenberg.txt")

	calculateCodeForEachChar(root)

	TraverseTree(root)

}

// TODO create cli api
func readFile(filename string) *os.File {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	return file
}
