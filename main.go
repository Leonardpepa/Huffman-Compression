package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
	"unicode/utf8"
)

type HuffmanTreeNode struct {
	left   *HuffmanTreeNode
	right  *HuffmanTreeNode
	weight uint64
	char   rune
	isLeaf bool
	code   string
}

func main() {
	fmt.Println("Hello, huffman!")

	//file, err := os.Open("gutenberg.txt")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//defer func(file *os.File) {
	//	err := file.Close()
	//	if err != nil {
	//
	//	}
	//}(file)
	//
	//reader := bufio.NewReader(file)
	//charFrequencies, err := calculateFrequencies(reader)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}

	charFrequencies := map[rune]uint64{
		'Z': 2,
		'K': 7,
		'M': 24,
		'C': 32,
		'U': 37,
		'D': 42,
		'L': 42,
		'E': 120,
	}

	treeNodes := createSliceFromMap(charFrequencies)

	root := buildHuffmanTree(treeNodes)

	table := make(map[rune]string)

	calculateCodeForEachChar(root, table)

	traverseTree(root)

	reader := bufio.NewReader(bytes.NewReader([]byte("EU")))
	compressData(root, table, reader)
}

func compressData(root *HuffmanTreeNode, table map[rune]string, reader *bufio.Reader) {
	encodedTree := encodeHuffmanHeaderInformation(root)

	output, err := os.Create("encoded.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func(output *os.File) {
		err := output.Close()
		if err != nil {
		}
	}(output)

	// number of characters in header
	_, err = output.WriteString(string(rune(utf8.RuneCountInString(encodedTree))))

	_, err = output.WriteString(encodedTree)
	if err != nil {
		log.Fatal(err)
	}

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		// TODO: store bits not strings
		_, err = output.WriteString(table[r])
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createLeafNode(char rune, freq uint64) *HuffmanTreeNode {
	return &HuffmanTreeNode{left: nil, right: nil, char: char, weight: freq, isLeaf: true}
}

func createHuffmanNode(a *HuffmanTreeNode, b *HuffmanTreeNode) *HuffmanTreeNode {
	var leftNode *HuffmanTreeNode
	var rightNode *HuffmanTreeNode
	if a.weight <= b.weight {
		leftNode = a
		rightNode = b
	} else {
		leftNode = b
		rightNode = a
	}
	return &HuffmanTreeNode{
		left:   leftNode,
		right:  rightNode,
		weight: leftNode.weight + rightNode.weight,
		isLeaf: false,
	}
}

func buildHuffmanTree(nodes []*HuffmanTreeNode) *HuffmanTreeNode {
	var a *HuffmanTreeNode
	var b *HuffmanTreeNode
	var root *HuffmanTreeNode
	sortHuffmanSlice(nodes)

	for len(nodes) > 1 {
		nodes, a = deleteAndReturn(nodes, 0)
		nodes, b = deleteAndReturn(nodes, 0)

		root = createHuffmanNode(a, b)

		nodes = append(nodes, root)
		sortHuffmanSlice(nodes)
	}
	return root
}

func sortHuffmanSlice(array []*HuffmanTreeNode) {
	slices.SortFunc(array, func(a, b *HuffmanTreeNode) int {
		return int(a.weight - b.weight)
	})
}

// create priority queue
func createSliceFromMap(charFrequencies map[rune]uint64) []*HuffmanTreeNode {
	treeNodes := make([]*HuffmanTreeNode, 0)

	for key, val := range charFrequencies {
		treeNodes = append(treeNodes, createLeafNode(key, val))
	}

	return treeNodes
}
func deleteAndReturn(slice []*HuffmanTreeNode, index int) ([]*HuffmanTreeNode, *HuffmanTreeNode) {
	item := slice[index]
	return append(slice[:index], slice[index+1:]...), item
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
	if node.isLeaf {
		node.code = string(c)
		table[node.char] = node.code
	}
	if node.right != nil {
		r := slices.Clone(c)
		r = append(r, '1')

		calculateCode(node.right, table, r)
	}
}

func traverseTree(root *HuffmanTreeNode) {
	if root.left != nil {
		traverseTree(root.left)
	}
	if root.isLeaf {
		fmt.Printf("char: %#v, freq: %d, code: %v, bits: %d\n", string(root.char), root.weight, root.code, len(root.code))
	}
	if root.right != nil {
		traverseTree(root.right)
	}
}

func calculateFrequencies(reader *bufio.Reader) (map[rune]uint64, error) {
	frequencies := make(map[rune]uint64)

	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		frequencies[char]++
	}

	return frequencies, nil
}

func encodeHuffmanHeaderInformation(node *HuffmanTreeNode) string {
	var builder strings.Builder
	// root node
	builder.WriteRune('0')
	recursiveHeaderEncoding(node, &builder)

	return builder.String()
}

func recursiveHeaderEncoding(node *HuffmanTreeNode, builder *strings.Builder) {
	if node == nil {
		log.Fatal("error nil pointer given for header encoding")
	}

	if node.left != nil {
		recursiveHeaderEncoding(node.left, builder)
	}
	if node.isLeaf {
		builder.WriteByte('1')
		builder.WriteRune(node.char)
	} else {
		builder.WriteByte('0')
	}

	if node.right != nil {
		recursiveHeaderEncoding(node.right, builder)
	}

}

func decodeString(node *HuffmanTreeNode, code string) (string, error) {
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
