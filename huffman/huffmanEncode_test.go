package huffman

//
//import (
//	"huffmanCompression/bitstream"
//	"testing"
//)
//
//func TestEncodeHuffmanHeaderInformation(t *testing.T) {
//	t.Run("should write the size followed by the tree in preorder traversal", func(t *testing.T) {
//		frequencies := map[rune]uint64{
//			'C': 32,
//			'D': 42,
//			'E': 120,
//			'K': 7,
//			'L': 42,
//			'M': 24,
//			'U': 37,
//			'Z': 2,
//		}
//		root := CreateHuffmanTreeFromFrequencies(frequencies)
//
//		writer := bitstream.CreateBitWriter()
//
//		size, err := encodeHuffmanHeaderInformation(root, PreOrder, &writer)
//
//		if err != nil {
//			t.Error(err)
//		}
//
//	})
//}
