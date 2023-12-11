package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
)

type HuffmanTree struct {
	left   *HuffmanTree
	right  *HuffmanTree
	weight uint64
	char   rune
	isLeaf bool
	code   string
}

func createTreeNode(char rune, freq uint64) *HuffmanTree {
	return &HuffmanTree{left: nil, right: nil, char: char, weight: freq, isLeaf: true}
}

func createHuffmanNode(a *HuffmanTree, b *HuffmanTree) *HuffmanTree {
	var leftNode *HuffmanTree
	var rightNode *HuffmanTree
	if a.weight <= b.weight {
		leftNode = a
		rightNode = b
	} else {
		leftNode = b
		rightNode = a
	}
	return &HuffmanTree{
		left:   leftNode,
		right:  rightNode,
		weight: leftNode.weight + rightNode.weight,
		isLeaf: false,
	}
}

func createTree(nodes []*HuffmanTree) *HuffmanTree {
	var a *HuffmanTree
	var b *HuffmanTree
	var root *HuffmanTree
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

func sortHuffmanSlice(array []*HuffmanTree) {
	slices.SortFunc(array, func(a, b *HuffmanTree) int {
		return int(a.weight - b.weight)
	})
}

func createSliceFromMap(charFrequencies map[rune]uint64) []*HuffmanTree {
	treeNodes := make([]*HuffmanTree, 0)

	for key, val := range charFrequencies {
		treeNodes = append(treeNodes, createTreeNode(key, val))
	}

	return treeNodes
}

func main() {
	fmt.Println("Hello, huffman!")

	file, err := os.Open("gutenberg.txt")
	if err != nil {
		log.Fatal(err)
	}

	charFrequencies, err := calculateFrequencies(bufio.NewReader(file))

	if err != nil {
		log.Fatal(err)
	}

	//charFrequencies := map[rune]uint64{
	//	'Z': 2,
	//	'K': 7,
	//	'M': 24,
	//	'C': 32,
	//	'U': 37,
	//	'D': 42,
	//	'L': 42,
	//	'E': 120,
	//}

	treeNodes := createSliceFromMap(charFrequencies)

	root := createTree(treeNodes)

	calculateCodeForEachChar(root)

	traverseTree(root)

}

func deleteElement(slice []*HuffmanTree, index int) ([]*HuffmanTree, *HuffmanTree) {
	item := slice[index]
	return append(slice[:index], slice[index+1:]...), item
}

func calculateCodeForEachChar(root *HuffmanTree) {
	var c []byte
	traverseWithCode(root, c)
}

func traverseWithCode(root *HuffmanTree, c []byte) {
	if root.left != nil {
		l := slices.Clone(c)
		l = append(l, byte('0'))

		traverseWithCode(root.left, l)
	}
	if root.isLeaf {
		root.code = string(c)
	}
	if root.right != nil {
		r := slices.Clone(c)
		r = append(r, byte('1'))

		traverseWithCode(root.right, r)
	}
}

func traverseTree(root *HuffmanTree) {
	if root.left != nil {
		traverseTree(root.left)
	}
	if root.isLeaf {
		s := fmt.Sprintf("char: %#v, freq: %d, code: %s\n", string(root.char), root.weight, root.code)
		fmt.Print(s)
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
