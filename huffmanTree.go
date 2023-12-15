package main

import (
	"container/heap"
	"fmt"
)

type HuffmanTreeNode struct {
	Left   *HuffmanTreeNode
	Right  *HuffmanTreeNode
	Weight uint64
	Char   rune
	IsLeaf bool
	Code   string
	index  int
}

type PriorityQueue []*HuffmanTreeNode

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Weight < pq[j].Weight
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
	item.Weight = weight
	heap.Fix(pq, item.index)
}

func CreateLeafNode(char rune, freq uint64) *HuffmanTreeNode {
	return &HuffmanTreeNode{Left: nil, Right: nil, Char: char, Weight: freq, IsLeaf: true}
}

func CreateHuffmanNode(a *HuffmanTreeNode, b *HuffmanTreeNode) *HuffmanTreeNode {
	var leftNode *HuffmanTreeNode
	var rightNode *HuffmanTreeNode
	if a.Weight <= b.Weight {
		leftNode = a
		rightNode = b
	} else {
		leftNode = b
		rightNode = a
	}
	return &HuffmanTreeNode{
		Left:   leftNode,
		Right:  rightNode,
		Weight: leftNode.Weight + rightNode.Weight,
		IsLeaf: false,
	}
}

func BuildHuffmanTree(pq PriorityQueue) *HuffmanTreeNode {
	var a *HuffmanTreeNode
	var b *HuffmanTreeNode
	var root *HuffmanTreeNode

	for pq.Len() > 1 {
		a = heap.Pop(&pq).(*HuffmanTreeNode)
		b = heap.Pop(&pq).(*HuffmanTreeNode)
		root = CreateHuffmanNode(a, b)

		heap.Push(&pq, root)
	}
	return root
}

// create priority queue
func CreatePriorityQueue(charFrequencies map[rune]uint64) PriorityQueue {
	pq := make(PriorityQueue, 0)

	i := 0
	for key, val := range charFrequencies {
		node := CreateLeafNode(key, val)
		node.index = i
		heap.Push(&pq, node)
		i++
	}

	heap.Init(&pq)

	return pq
}

func TraverseTree(root *HuffmanTreeNode) {
	if root.Left != nil {
		TraverseTree(root.Left)
	}
	if root != nil && root.IsLeaf {
		fmt.Printf("Char: %v, %#v, freq: %d, Code: %v, bits: %d\n", root.Char, string(root.Char), root.Weight, root.Code, len(root.Code))
	}
	if root.Right != nil {
		TraverseTree(root.Right)
	}
}
