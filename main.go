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

	input := "hamlet"
	file := readFile("input/" + input + ".txt")

	root, err := getHuffmanTreeFromFile(file)

	if err != nil {
		log.Fatal(err)
	}

	calculateCodeForEachChar(root)

	TraverseTree(root)

	//////without header
	//encodedData, err := createBits(file, table)
	//if err != nil {
	//	return
	//}
	//
	//output := "output/" + input + ".hf"
	//_ = os.WriteFile(output, encodedData, 0666)
	//
	//text, err := getDecodedText(root, output)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println(*text)

	information, size, err := encodeHuffmanHeaderInformation(root, PostOrder)
	if err != nil {
		log.Fatal(err)
	}

	reader := CreateBitReader(information)

	createTreeFromHeader(&reader, size)
}

func createTreeFromHeader(reader *BitReader, size int) {
	count := 0
	for count < size && reader.HasNext() {
		bit := reader.Read()
		switch bit {
		case true:
			c := reader.ReadChar()
			fmt.Print(1)
			fmt.Printf("%c", c)
			count++
		case false:
			fmt.Print(0)
		}
		count++
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
