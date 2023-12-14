package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

func main() {
	var charFrequencies map[rune]uint64

	//file := readFile("gutenberg.txt")
	//
	//defer func(file *os.File) {
	//	err := file.Close()
	//	if err != nil {
	//
	//	}
	//}(file)

	//charFrequencies, err := calculateFrequencies(bufio.NewReader(file))
	//if err != nil {
	//	log.Fatal(err)
	//}

	charFrequencies = map[rune]uint64{
		'Z': 2,
		'K': 7,
		'M': 24,
		'C': 32,
		'U': 37,
		'D': 42,
		'L': 42,
		'E': 120,
	}

	priorityQueue := createPriorityQueue(charFrequencies)

	root := buildHuffmanTree(priorityQueue)

	table := make(map[rune]string)

	calculateCodeForEachChar(root, table)

	traverseTree(root)
}

// TODO take bits and format it correctly in strings
func formatBitString(byteArray []byte) string {
	strToBeDecoded := fmt.Sprintf("%b", byteArray)
	strToBeDecoded = strings.TrimFunc(strings.Join(strings.Fields(strToBeDecoded), ""), func(r rune) bool {
		return r == '[' || r == ']'
	})
	return strToBeDecoded
}

// TODO fix this
func compressData(table map[rune]string, reader *bufio.Reader) []byte {

	byteArray := make([]byte, 0)
	x := byte(0)

	for {
		c, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		count := 0

		for _, bitt := range table[c] {

			if bitt == '1' {
				x = x<<1 | 1
			} else if bitt == '0' {
				x = x << 1
			}
			count++

			if count == 8 || count == len(table[c]) {
				byteArray = append(byteArray, x)
				count = 0
				x = byte(0)
			}
		}
	}
	return byteArray
}

// TODO create cli api
func readFile(filename string) *os.File {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	return file
}

func calculateCodeForEachChar(node *HuffmanTreeNode, table map[rune]string) {
	var c []byte
	calculateCode(node, table, c)
}

func calculateCode(node *HuffmanTreeNode, table map[rune]string, c []byte) {
	if node.left != nil {
		l := slices.Clone(c)
		l = append(l, '0')

		calculateCode(node.left, table, l)
	}
	if node != nil && node.isLeaf {
		node.code = string(c)
		table[node.char] = node.code
	}
	if node.right != nil {
		r := slices.Clone(c)
		r = append(r, '1')

		calculateCode(node.right, table, r)
	}
}

func encodeHuffmanHeaderInformation(node *HuffmanTreeNode) (string, error) {
	// root node
	var strBuilder strings.Builder
	strBuilder.WriteRune('0')

	err := recursiveHeaderEncoding(node, &strBuilder)
	return strBuilder.String(), err
}

func recursiveHeaderEncoding(node *HuffmanTreeNode, builder *strings.Builder) error {
	var err error
	if node == nil {
		return fmt.Errorf("error nil pointer given for header encoding")
	}

	if node.left != nil {
		err = recursiveHeaderEncoding(node.left, builder)
	}
	if node.isLeaf {
		builder.WriteRune('1')
		builder.WriteRune(node.char)
	} else {
		builder.WriteRune('0')
	}

	if node.right != nil {
		err = recursiveHeaderEncoding(node.right, builder)
	}
	return err
}

func decodeText(node *HuffmanTreeNode, code string) (string, error) {
	var strBuilder strings.Builder
	temp := node

	for _, val := range code {
		switch val {
		case '0':
			temp = temp.left
		case '1':
			temp = temp.right
		default:
			return "", fmt.Errorf("something went wrong, input != 0 || 1")
		}

		if temp == nil {
			return "", fmt.Errorf("nil reference")
		}
		if temp.isLeaf {
			strBuilder.WriteRune(temp.char)
			temp = node
			continue
		}

	}
	return strBuilder.String(), nil
}
