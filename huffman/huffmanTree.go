package huffman

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"log"
	"os"
	"unicode/utf8"
)

type TreeNode struct {
	Left   *TreeNode
	Right  *TreeNode
	Parent *TreeNode
	Weight uint64
	Char   rune
	IsLeaf bool
	Code   string
	index  int
}

type PriorityQueue []*TreeNode

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
	item := x.(*TreeNode)
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
func (pq *PriorityQueue) update(item *TreeNode, weight uint64) {
	item.Weight = weight
	heap.Fix(pq, item.index)
}

func CreateLeafNode(char rune, freq uint64) *TreeNode {
	return &TreeNode{Left: nil, Right: nil, Parent: nil, Char: char, Weight: freq, IsLeaf: true}
}

func CreateHuffmanNode(a *TreeNode, b *TreeNode) *TreeNode {
	var leftNode *TreeNode
	var rightNode *TreeNode
	if a.Weight < b.Weight {
		leftNode = a
		rightNode = b
	} else {
		leftNode = b
		rightNode = a
	}
	parent := &TreeNode{
		Left:   leftNode,
		Right:  rightNode,
		Weight: leftNode.Weight + rightNode.Weight,
		IsLeaf: false,
	}
	leftNode.Parent = parent
	rightNode.Parent = parent
	return parent
}

func BuildHuffmanTree(pq PriorityQueue) *TreeNode {
	var a *TreeNode
	var b *TreeNode
	var root *TreeNode

	for pq.Len() > 1 {
		a = heap.Pop(&pq).(*TreeNode)
		b = heap.Pop(&pq).(*TreeNode)
		root = CreateHuffmanNode(a, b)

		heap.Push(&pq, root)
	}
	return root
}

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

func TraverseTreeInorder(root *TreeNode) {
	if root.Left != nil {
		TraverseTreeInorder(root.Left)
	}

	if root != nil && root.IsLeaf {
		p := make([]byte, utf8.RuneLen(root.Char))
		utf8.EncodeRune(p, root.Char)
		fmt.Printf("bytes: %b, Char: %v, %#v, freq: %d, Code: %v, bits: %d\n", p, root.Char, string(root.Char), root.Weight, root.Code, len(root.Code))
	}
	if root.Right != nil {
		TraverseTreeInorder(root.Right)
	}
}

func CalculateCodeForEachChar(node *TreeNode) map[rune]string {
	table := make(map[rune]string)
	calculateCode(node, table, "")
	return table
}

func calculateCode(node *TreeNode, table map[rune]string, c string) {
	if node.Left != nil {
		calculateCode(node.Left, table, c+"0")
	}
	if node != nil && node.IsLeaf {
		node.Code = c
		table[node.Char] = node.Code
	}
	if node.Right != nil {
		calculateCode(node.Right, table, c+"1")
	}
}

func CreateHuffmanTreeFromFrequencies(charFrequencies map[rune]uint64) *TreeNode {

	priorityQueue := CreatePriorityQueue(charFrequencies)

	root := BuildHuffmanTree(priorityQueue)

	return root
}

func CreateHuffmanTreeFromFile(file *os.File) (*TreeNode, error) {
	_, err := file.Seek(0, io.SeekStart)

	if err != nil {
		return nil, err
	}

	log.Println("Calculating frequencies...")
	charFrequencies, err := calculateFrequencies(bufio.NewReader(file))

	if err != nil {
		return nil, err
	}

	AddPseudoEOF(charFrequencies)

	log.Println("Create priority queue...")
	priorityQueue := CreatePriorityQueue(charFrequencies)

	log.Println("Build the tree...")
	root := BuildHuffmanTree(priorityQueue)

	return root, nil
}
