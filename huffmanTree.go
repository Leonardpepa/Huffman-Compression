package main

import (
	"fmt"
	"slices"
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

func buildHuffmanTree(nodes []*HuffmanTreeNode) *HuffmanTreeNode {
	var a *HuffmanTreeNode
	var b *HuffmanTreeNode
	var root *HuffmanTreeNode
	sortHuffmanSlice(nodes)

	for len(nodes) > 1 {
		a = deleteAndReturn(&nodes, 0)
		b = deleteAndReturn(&nodes, 0)

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
func deleteAndReturn(slice *[]*HuffmanTreeNode, index int) *HuffmanTreeNode {
	item := (*slice)[index]
	*slice = append((*slice)[:index], (*slice)[index+1:]...)
	return item
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
