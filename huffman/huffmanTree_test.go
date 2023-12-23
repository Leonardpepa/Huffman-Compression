package huffman

import (
	"container/heap"
	"maps"
	"slices"
	"testing"
)

// https://opendsa-server.cs.vt.edu/ODSA/Books/CS3/html/Huffman.html
func TestCreatePriorityQueue(t *testing.T) {
	t.Run("create priority queue from frequency table", func(t *testing.T) {
		frequencies := map[rune]uint64{
			'M': 24,
			'C': 32,
			'K': 7,
			'D': 42,
			'E': 120,
			'L': 42,
			'U': 37,
			'Z': 2,
		}

		pq := CreatePriorityQueue(frequencies)

		prev := uint64(0)
		for pq.Len() > 0 {
			item := heap.Pop(&pq).(*TreeNode)

			if item.Weight < prev {
				t.Errorf("wrong item popped. last was %d but got %d", prev, item.Weight)
			}

			prev = item.Weight
		}

		heap.Push(&pq, CreateLeafNode('P', 163))
		heap.Push(&pq, CreateLeafNode('0', 1))

		item := heap.Pop(&pq).(*TreeNode)

		if item.Weight != 1 {
			t.Errorf("wrong item popped. expected 1, got %d", item.Weight)
		}

		if pq.Len() != 1 {
			t.Errorf("wron priority queue size expected 1. got %d", pq.Len())
		}

	})
}

func TestBuildHuffmanTree(t *testing.T) {
	t.Run("return a huffman tree", func(t *testing.T) {
		frequencies := map[rune]uint64{
			'Z': 2,
			'K': 7,
			'M': 24,
			'C': 32,
			'D': 42,
			'L': 42,
			'U': 37,
			'E': 120,
		}

		root := CreateHuffmanTreeFromFrequencies(frequencies)

		charsInorder := make([]rune, 0)

		getInorderChars(root, &charsInorder)

		expected := []rune{'E', 'U', 'D', 'L', 'C', 'Z', 'K', 'M'}

		if !slices.EqualFunc(expected, charsInorder, func(r rune, r2 rune) bool { return r == r2 }) {
			t.Errorf("something went wrong when building the tree, expected %v, got %v", expected, charsInorder)
		}
	})
}

func TestCalculateBitCodes(t *testing.T) {
	t.Run("create bit codes for each character", func(t *testing.T) {
		frequencies := map[rune]uint64{
			'Z': 2,
			'K': 7,
			'M': 24,
			'C': 32,
			'D': 42,
			'L': 42,
			'U': 37,
			'E': 120,
		}

		root := CreateHuffmanTreeFromFrequencies(frequencies)

		table := CalculateBitCodes(root)

		expected := map[rune]string{
			'M': "11111",
			'C': "1110",
			'K': "111101",
			'D': "101",
			'E': "0",
			'L': "110",
			'U': "100",
			'Z': "111100",
		}

		if !maps.Equal(expected, table) {
			t.Errorf("something went wrong when calculating the bit codes expected %v, got %v", expected, table)
		}

	})
}
