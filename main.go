package main

import (
	"fmt"
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
	//AddPseudoEOF(charFrequencies)
	//
	//root := getHuffmanTreeFromFrequencies(charFrequencies)

	input := "gogo"
	file := readFile("input/" + input + ".txt")

	root, err := getHuffmanTreeFromFile(file)

	if err != nil {
		log.Fatal(err)
	}

	table := calculateCodeForEachChar(root)

	TraverseTree(root)

	//without header
	encodedData, err := createBits(file, table)
	if err != nil {
		return
	}

	output := "output/" + input + ".hf"
	_ = os.WriteFile(output, encodedData, 0666)

	information, size, err := encodeHuffmanHeaderInformation(root, PreOrder)
	if err != nil {
		log.Fatal(err)
	}

	reader := CreateBitReader(information)

	rootFromHeader := CreateTreeFromHeader(&reader, size)
	text, err := getDecodedText(rootFromHeader, output)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*text)

}

// TODO create cli api
func readFile(filename string) *os.File {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	return file
}
