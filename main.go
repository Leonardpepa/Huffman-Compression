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
	//
	//charFrequencies, err := CalculateFrequencies(bufio.NewReader(file))
	//AddPseudoEOF(charFrequencies)
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

	priorityQueue := CreatePriorityQueue(charFrequencies)

	root := BuildHuffmanTree(priorityQueue)

	table := make(map[rune]string)

	calculateCodeForEachChar(root, table)

	TraverseTree(root)

	information, err := encodeHuffmanHeaderInformation(root)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(information)

	newRoot := constructTreeFromHeader(information)
	TraverseTree(newRoot)
}

func constructTreeFromHeader(information string) *HuffmanTreeNode {
	newRoot := &HuffmanTreeNode{}
	var prev *HuffmanTreeNode
	var curr *HuffmanTreeNode

	hold := ""
	for _, c := range information {
		hold += string(c)
		switch c {
		case '0':
			if curr == nil {
				newRoot = &HuffmanTreeNode{}
				curr = newRoot
				continue
			}

			if curr.Left == nil {
				curr.Left = &HuffmanTreeNode{Parent: curr}
				prev = curr
				curr = curr.Left
			} else if curr.Right == nil {
				curr.Right = &HuffmanTreeNode{Parent: curr}
				prev = curr
				curr = curr.Right
			}

		case '1':
			if curr.Left == nil {
				curr.Left = &HuffmanTreeNode{Parent: curr}
				prev = curr
				curr = curr.Left
			} else if curr.Right == nil {
				curr.Right = &HuffmanTreeNode{Parent: curr}
				prev = curr.Parent
				curr = curr.Right
			}
		default:
			curr.IsLeaf = true
			curr.Char = c
			curr = prev
		}
	}
	return newRoot
}

func getDecodedText(root *HuffmanTreeNode, filename string) string {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	res := fmt.Sprintf("%08b", fileBytes)
	res = strings.TrimFunc(res, func(r rune) bool {
		return r == '[' || r == ']'
	})
	res = strings.Join(strings.Fields(res), "")
	text, err := decodeText(root, res)
	if err != nil {
		log.Fatal(err)
	}
	return text
}

func createBits(file *os.File, table map[rune]string) []byte {
	compressedString := getCompressedDataAsString(file, table)
	arrayOfBytes := make([]byte, 0)
	count := 0
	compressedLength := len(compressedString)
	x := byte(0)

	for _, value := range compressedString {
		switch value {
		case '0':
			x <<= 1
		case '1':
			x = x<<1 | 1
		}
		count++

		if count == 8 || count == compressedLength {
			arrayOfBytes = append(arrayOfBytes, x)
			count = 0
			x = byte(0)
		}
	}
	return arrayOfBytes
}

func getCompressedDataAsString(file *os.File, table map[rune]string) string {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}
	var builder strings.Builder
	reader := bufio.NewReader(file)

	for {
		r, _, err := reader.ReadRune()

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		builder.WriteString(table[r])
	}

	// write the Pseudo EOF in the end of the characters
	builder.WriteString(table[PseudoEOF])
	return builder.String()
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
	if node.Left != nil {
		l := slices.Clone(c)
		l = append(l, '0')

		calculateCode(node.Left, table, l)
	}
	if node != nil && node.IsLeaf {
		node.Code = string(c)
		table[node.Char] = node.Code
	}
	if node.Right != nil {
		r := slices.Clone(c)
		r = append(r, '1')

		calculateCode(node.Right, table, r)
	}
}

func encodeHuffmanHeaderInformation(node *HuffmanTreeNode) (string, error) {
	// root node
	var strBuilder strings.Builder
	err := recursiveHeaderEncoding(node, &strBuilder)
	if err != nil {
		return "", err
	}
	return strBuilder.String(), nil
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
			return "", fmt.Errorf("something went wrong, input != 0 || 1")
		}

		if temp == nil {
			return "", fmt.Errorf("nil reference")
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
