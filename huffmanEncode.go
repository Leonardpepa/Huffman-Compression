package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func encodeHuffmanHeaderInformation(node *HuffmanTreeNode) ([]byte, error) {
	// root node
	writer := CreateBitWriter()
	err := recursiveHeaderEncoding(node, &writer)
	// stop here
	writer.writeRune(PseudoEOF)
	writer.WriteBytes()

	if err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

func recursiveHeaderEncoding(node *HuffmanTreeNode, writer *BitWriter) error {
	var err error
	if node == nil {
		return fmt.Errorf("error nil pointer given for header encoding")
	}

	if node.IsLeaf {
		writer.writeBitFromBool(true)
		writer.writeRune(node.Char)
	} else {
		writer.writeBitFromBool(false)
	}

	if node.Left != nil {
		err = recursiveHeaderEncoding(node.Left, writer)
	}

	if node.Right != nil {
		err = recursiveHeaderEncoding(node.Right, writer)
	}

	return err
}

func createBits(file *os.File, table map[rune]string) ([]byte, error) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)
	bitWriter := CreateBitWriter()

	for {
		r, _, err := reader.ReadRune()

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		for _, value := range table[r] {
			bitWriter.writeBitFromChar(value)
		}
	}

	// Pseudo EOF
	for _, value := range table[PseudoEOF] {
		bitWriter.writeBitFromChar(value)
	}

	// finish the writing
	bitWriter.WriteBytes()

	return bitWriter.Bytes(), nil
}
