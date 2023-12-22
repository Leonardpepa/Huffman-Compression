package huffman

import (
	"fmt"
	"huffmanCompression/bitstream"
	"io"
	"log"
	"os"
	"strings"
)

func Decode(file *os.File, output string) error {
	log.Println("Reading encoded data... ")
	data, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	bitReader := bitstream.CreateBitReader(data)
	sizeRune, err := bitReader.ReadChar()
	size := int(sizeRune)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Constructing the huffman tree from the header... ")
	root, err := createTreeFromHeader(&bitReader, size)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Decoding the data... ")
	text, err := getDecodedText(root, &bitReader)
	if err != nil {
		return err
	}

	err = os.WriteFile(output, []byte(*text), 0666)

	if err != nil {
		return err
	}

	log.Println("Decoding finished. File: ", output)
	return nil
}

func getDecodedText(root *TreeNode, bitReader *bitstream.Reader) (*string, error) {

	text, err := decodeText(root, bitReader)

	if err != nil {
		return nil, err
	}
	return &text, nil
}

func decodeText(node *TreeNode, bitReader *bitstream.Reader) (string, error) {
	var strBuilder strings.Builder
	temp := node

	for bitReader.HasNext() {
		bit, err := bitReader.Read()
		if err != nil {
			return "", err
		}

		switch bit {
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

func createTreeFromHeader(reader *bitstream.Reader, size int) (*TreeNode, error) {
	var root *TreeNode
	var current *TreeNode

	count := 0
	for count < size && reader.HasNext() {
		bit, err := reader.Read()
		if err != nil {
			return nil, err
		}
		switch bit {
		case true:
			count += 2
			c, err := reader.ReadChar()
			if err != nil {
				return nil, err
			}
			current = backtrack(current)
			createLeafNode(current, c)

		case false:
			count++
			if current == nil {
				root = &TreeNode{}
				current = root
				continue
			}

			current = backtrack(current)
			current = createInternalNode(current)
		}
	}

	return root, nil
}

func createInternalNode(current *TreeNode) *TreeNode {
	if current.Left == nil {
		temp := &TreeNode{Parent: current}
		current.Left = temp
		current = current.Left
	} else if current.Right == nil {
		temp := &TreeNode{Parent: current}
		current.Right = temp
		current = current.Right
	} else {
		log.Fatal("error while creating internal node")
	}

	return current
}

func createLeafNode(current *TreeNode, c rune) {
	if current.Left == nil {
		temp := &TreeNode{Parent: current, IsLeaf: true, Char: c}
		current.Left = temp
	} else if current.Right == nil {
		temp := &TreeNode{Parent: current, IsLeaf: true, Char: c}
		current.Right = temp
	}
}

func backtrack(current *TreeNode) *TreeNode {
	for current.Left != nil && current.Right != nil {
		if current.Parent == nil {
			log.Fatal("Error while decoding tree, nil pointer")
		}
		current = current.Parent
	}
	return current
}
