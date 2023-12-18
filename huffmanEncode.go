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
	data, err := getCompressedDataAsString(file, table)
	if err != nil {
		return nil, err
	}

	compressedString := *data
	bitWriter := CreateBitWriter()

	for _, value := range compressedString {
		bitWriter.writeBitFromChar(value)
	}

	if bitWriter.HasRemainingBits() {
		bitWriter.WriteRemainingBitsWithPadding()
	}
	return bitWriter.Bytes(), nil
}

func getCompressedDataAsString(file *os.File, table map[rune]string) (*string, error) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	var builder strings.Builder
	reader := bufio.NewReader(file)

	for {
		r, _, err := reader.ReadRune()

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		builder.WriteString(table[r])
	}

	// write the Pseudo EOF in the end of the characters
	builder.WriteString(table[PseudoEOF])
	data := builder.String()
	return &data, nil
}
