package main

import (
	"container/heap"
	"fmt"
	"sort"
)

type HuffmanTreeNode struct {
	left   *HuffmanTreeNode
	right  *HuffmanTreeNode
	weight uint64
	char   rune
	isLeaf bool
	code   string
	index  int
}

type PriorityQueue []*HuffmanTreeNode

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].weight < pq[j].weight
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*HuffmanTreeNode)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *HuffmanTreeNode, weight uint64) {
	item.weight = weight
	heap.Fix(pq, item.index)
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
	sort.Slice(array, func(i, j int) bool {
		return array[i].weight < array[j].weight
	})
}

// create priority queue
func createPriorityQueue(charFrequencies map[rune]uint64) PriorityQueue {
	pq := make(PriorityQueue, 0)

	i := 0
	for key, val := range charFrequencies {
		node := createLeafNode(key, val)
		node.index = i
		heap.Push(&pq, node)
		i++
	}

	heap.Init(&pq)

	return pq
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
