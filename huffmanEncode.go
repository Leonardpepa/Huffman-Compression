package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"unicode/utf8"
)

const (
	PreOrder = iota
	PostOrder
	InOrder
)

func Encode(file *os.File, output string) error {
	log.Println("Creating huffman tree... ")
	root, err := CreateHuffmanTreeFromFile(file)

	if err != nil {
		return err
	}

	log.Println("Calculating variable length codes... ")
	table := calculateCodeForEachChar(root)

	bitWriter := CreateBitWriter()

	log.Println("Encoding the tree in the header... ")
	size, err := encodeHuffmanHeaderInformation(root, PreOrder, &bitWriter)

	if err != nil {
		return err
	}

	log.Println("Encoding the data bit by bit... ")
	err = createBits(file, &bitWriter, table)

	if err != nil {
		return err
	}

	bitWriter.WriteBytes()

	sizeRune := rune(size)

	b := make([]byte, utf8.RuneLen(sizeRune))
	utf8.EncodeRune(b, sizeRune)

	data := make([]byte, 0)

	data = append(data, b...)
	data = append(data, bitWriter.Bytes()...)

	err = os.WriteFile(output, data, 0666)

	if err != nil {
		return err
	}

	log.Println("Encoding finished. File: ", output)
	return nil
}

func encodeHuffmanHeaderInformation(node *HuffmanTreeNode, encodeType int, writer *BitWriter) (int, error) {
	var err error
	// root node
	count := 0
	switch encodeType {
	case PreOrder:
		err = recursivePreOrderHeaderEncoding(node, writer, &count)
	case PostOrder:
		err = recursivePostOrderHeaderEncoding(node, writer, &count)
	case InOrder:
		err = recursiveInOrderHeaderEncoding(node, writer, &count)
	}
	//// stop here
	//writer.WriteBytes()

	if err != nil {
		return 0, err
	}
	return count, nil
}

func recursivePreOrderHeaderEncoding(node *HuffmanTreeNode, writer *BitWriter, count *int) error {
	var err error
	if node == nil {
		return fmt.Errorf("error nil pointer given for header encoding")
	}

	if node.IsLeaf {
		writer.writeBitFromBool(true)
		writer.writeRune(node.Char)
		*count = *count + 2
	} else {
		writer.writeBitFromBool(false)
		*count = *count + 1
	}

	if node.Left != nil {
		err = recursivePreOrderHeaderEncoding(node.Left, writer, count)
	}

	if node.Right != nil {
		err = recursivePreOrderHeaderEncoding(node.Right, writer, count)
	}

	return err
}

func recursivePostOrderHeaderEncoding(node *HuffmanTreeNode, writer *BitWriter, count *int) error {
	var err error
	if node == nil {
		return fmt.Errorf("error nil pointer given for header encoding")
	}

	if node.Left != nil {
		err = recursivePostOrderHeaderEncoding(node.Left, writer, count)
	}

	if node.Right != nil {
		err = recursivePostOrderHeaderEncoding(node.Right, writer, count)
	}

	if node.IsLeaf {
		writer.writeBitFromBool(true)
		writer.writeRune(node.Char)
		*count = *count + 2
	} else {
		writer.writeBitFromBool(false)
		*count = *count + 1
	}

	return err
}

func recursiveInOrderHeaderEncoding(node *HuffmanTreeNode, writer *BitWriter, count *int) error {
	var err error
	if node == nil {
		return fmt.Errorf("error nil pointer given for header encoding")
	}

	if node.Left != nil {
		err = recursiveInOrderHeaderEncoding(node.Left, writer, count)
	}

	if node.IsLeaf {
		writer.writeBitFromBool(true)
		writer.writeRune(node.Char)
		*count = *count + 2
	} else {
		writer.writeBitFromBool(false)
		*count = *count + 1
	}
	if node.Right != nil {
		err = recursiveInOrderHeaderEncoding(node.Right, writer, count)
	}

	return err
}

func createBits(file *os.File, bitWriter *BitWriter, table map[rune]string) error {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)

	for {
		r, _, err := reader.ReadRune()

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		for _, value := range table[r] {
			bitWriter.writeBitFromChar(value)
		}
	}

	// Pseudo EOF
	for _, value := range table[PseudoEOF] {
		bitWriter.writeBitFromChar(value)
	}

	//// finish the writing
	//bitWriter.WriteBytes()

	return nil
}
