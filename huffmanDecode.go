package main

import (
	"fmt"
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
		switch bitReader.GetNextBit() {
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
