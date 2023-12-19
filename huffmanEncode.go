package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	PreOrder = iota
	PostOrder
	InOrder
)

func encodeHuffmanHeaderInformation(node *HuffmanTreeNode, encodeType int) ([]byte, int, error) {
	var err error
	// root node
	writer := CreateBitWriter()
	count := 0
	switch encodeType {
	case PreOrder:
		err = recursivePreOrderHeaderEncoding(node, &writer, &count)
	case PostOrder:
		err = recursivePostOrderHeaderEncoding(node, &writer, &count)
	case InOrder:
		err = recursiveInOrderHeaderEncoding(node, &writer, &count)
	}
	// stop here
	writer.WriteBytes()

	if err != nil {
		return nil, 0, err
	}
	return writer.Bytes(), count, nil
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
