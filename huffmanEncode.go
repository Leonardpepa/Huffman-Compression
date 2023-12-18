package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func encodeHuffmanHeaderInformation(node *HuffmanTreeNode) (*string, error) {
	// root node
	var strBuilder strings.Builder
	err := recursiveHeaderEncoding(node, &strBuilder)
	if err != nil {
		return nil, err
	}
	data := strBuilder.String()
	return &data, nil
}

func recursiveHeaderEncoding(node *HuffmanTreeNode, builder *strings.Builder) error {
	var err error
	if node == nil {
		return fmt.Errorf("error nil pointer given for header encoding")
	}

	if node.IsLeaf {
		builder.WriteRune('1')
		builder.WriteRune(node.Char)
	} else {
		builder.WriteRune('0')
	}

	if node.Left != nil {
		err = recursiveHeaderEncoding(node.Left, builder)
	}

	if node.Right != nil {
		err = recursiveHeaderEncoding(node.Right, builder)
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
