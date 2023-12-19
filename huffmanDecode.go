package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func getDecodedText(root *HuffmanTreeNode, filename string) (*string, error) {
	fileBytes, err := os.ReadFile(filename)

	bitReader := CreateBitReader(fileBytes)

	text, err := decodeText(root, &bitReader)

	if err != nil {
		return nil, err
	}
	return &text, nil
}

func decodeText(node *HuffmanTreeNode, bitReader *BitReader) (string, error) {
	var strBuilder strings.Builder
	temp := node

	for bitReader.HasNext() {
		switch bitReader.Read() {
		case false:
			temp = temp.Left
		case true:
			temp = temp.Right
		default:
			return "", fmt.Errorf("something went wrong, input != (0 || 1)")
		}

		if temp == nil {
			return "", fmt.Errorf("nil reference pointer node when decoding text")
		}

		if temp.IsLeaf {
			// stop when you encounter Pseudo EOF
			if temp.Char == PseudoEOF {
				break
			}
			strBuilder.WriteRune(temp.Char)
			temp = node
		}

	}
	return strBuilder.String(), nil
}

func CreateTreeFromHeader(reader *BitReader, size int) *HuffmanTreeNode {
	var root *HuffmanTreeNode
	var current *HuffmanTreeNode

	count := 0
	for count < size && reader.HasNext() {
		bit := reader.Read()
		switch bit {
		case true:
			count += 2
			c := reader.ReadChar()
			current = backtrack(current)
			createLeafNode(current, c)

		case false:
			count++
			if current == nil {
				root = &HuffmanTreeNode{}
				current = root
				continue
			}

			current = backtrack(current)
			current = createInternalNode(current)
		}
	}

	return root
}

func createInternalNode(current *HuffmanTreeNode) *HuffmanTreeNode {
	if current.Left == nil {
		temp := &HuffmanTreeNode{Parent: current}
		current.Left = temp
		current = current.Left
	} else if current.Right == nil {
		temp := &HuffmanTreeNode{Parent: current}
		current.Right = temp
		current = current.Right
	} else {
		log.Fatal("error while creating internal node")
	}

	return current
}

func createLeafNode(current *HuffmanTreeNode, c rune) {
	if current.Left == nil {
		temp := &HuffmanTreeNode{Parent: current, IsLeaf: true, Char: c}
		current.Left = temp
	} else if current.Right == nil {
		temp := &HuffmanTreeNode{Parent: current, IsLeaf: true, Char: c}
		current.Right = temp
	}
}

func backtrack(current *HuffmanTreeNode) *HuffmanTreeNode {
	for current.Left != nil && current.Right != nil {
		if current.Parent == nil {
			log.Fatal("Error while decoding tree, nil pointer")
		}
		current = current.Parent
	}
	return current
}
