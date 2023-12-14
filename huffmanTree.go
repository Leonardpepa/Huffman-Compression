package main

import (
	"container/heap"
	"fmt"
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

func buildHuffmanTree(pq PriorityQueue) *HuffmanTreeNode {
	var a *HuffmanTreeNode
	var b *HuffmanTreeNode
	var root *HuffmanTreeNode

	for pq.Len() > 1 {
		a = heap.Pop(&pq).(*HuffmanTreeNode)
		b = heap.Pop(&pq).(*HuffmanTreeNode)
		root = createHuffmanNode(a, b)

		heap.Push(&pq, root)
	}
	return root
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
