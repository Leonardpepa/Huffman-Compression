package main

import (
	"fmt"
	"os"
	"strings"
)

func getDecodedText(root *HuffmanTreeNode, filename string) (*string, error) {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	res := fmt.Sprintf("%08b", fileBytes)
	res = strings.TrimFunc(res, func(r rune) bool {
		return r == '[' || r == ']'
	})
	res = strings.Join(strings.Fields(res), "")
	text, err := decodeText(root, res)

	if err != nil {
		return nil, err
	}
	return &text, nil
}

func decodeText(node *HuffmanTreeNode, code string) (string, error) {
	var strBuilder strings.Builder
	temp := node

	for _, val := range code {
		switch val {
		case '0':
			temp = temp.Left
		case '1':
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
