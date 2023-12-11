package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"slices"
	"strings"
)

type HuffmanTreeNode struct {
	left   *HuffmanTreeNode
	right  *HuffmanTreeNode
	weight uint64
	char   rune
	isLeaf bool
	code   string
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

func buildTree(nodes []*HuffmanTreeNode) *HuffmanTreeNode {
	var a *HuffmanTreeNode
	var b *HuffmanTreeNode
	var root *HuffmanTreeNode
	sortHuffmanSlice(nodes)

	for len(nodes) > 1 {
		nodes, a = deleteElement(nodes, 0)
		nodes, b = deleteElement(nodes, 0)

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

func createSliceFromMap(charFrequencies map[rune]uint64) []*HuffmanTreeNode {
	treeNodes := make([]*HuffmanTreeNode, 0)

	for key, val := range charFrequencies {
		treeNodes = append(treeNodes, createLeafNode(key, val))
	}

	return treeNodes
}

func main() {
	fmt.Println("Hello, huffman!")

	//file, err := os.Open("gutenberg.txt")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//charFrequencies, err := calculateFrequencies(bufio.NewReader(file))
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

	root := buildTree(treeNodes)
	calculateCodeForEachChar(root)

	traverseTree(root)
	decoded, err := decodeString(root, "1111010100")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v", decoded)
}

func deleteElement(slice []*HuffmanTreeNode, index int) ([]*HuffmanTreeNode, *HuffmanTreeNode) {
	item := slice[index]
	return append(slice[:index], slice[index+1:]...), item
}

func calculateCodeForEachChar(node *HuffmanTreeNode) {
	var c []byte
	traverseWithCode(node, c)
}

func traverseWithCode(node *HuffmanTreeNode, c []byte) {
	if node.left != nil {
		l := slices.Clone(c)
		l = append(l, '0')

		traverseWithCode(node.left, l)
	}
	if node.isLeaf {
		node.code = string(c)
	}
	if node.right != nil {
		r := slices.Clone(c)
		r = append(r, '1')

		traverseWithCode(node.right, r)
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

		if temp.isLeaf {
			strBuilder.WriteRune(temp.char)
			temp = node
			continue
		}

	}

	return strBuilder.String(), nil
}
