package huffman

import (
	"huffmanCompression/bitstream"
	"slices"
	"testing"
)

func TestDecodedText(t *testing.T) {
	t.Run("should decode the the bits", func(t *testing.T) {
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

		AddPseudoEOF(frequencies)

		root := CreateHuffmanTreeFromFrequencies(frequencies)

		table := CalculateBitCodes(root)
		writer := bitstream.CreateBitWriter()

		err := encodeBitCode(table, 'K', &writer)
		if err != nil {
			t.Error(err)
		}

		err = encodeBitCode(table, 'U', &writer)
		if err != nil {
			t.Error(err)
		}

		err = encodeBitCode(table, 'L', &writer)
		if err != nil {
			t.Error(err)
		}

		err = encodeBitCode(table, PseudoEOF, &writer)
		if err != nil {
			t.Error(err)
		}
		writer.Flush()

		reader := bitstream.CreateBitReader(writer.Bytes())

		text, err := getDecodedText(root, &reader)

		if err != nil {
			t.Error(err)
		}

		if *text != "KUL" {
			t.Errorf("when decoding text expected KUL, got %s", *text)
		}
	})
}

func TestCreateTreeFromHeader(t *testing.T) {
	t.Run("should create the tree by reading the header information", func(t *testing.T) {
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

		expected := make([]rune, 0)

		getInorderChars(root, &expected)

		writer := bitstream.CreateBitWriter()

		size, err := encodeHuffmanHeaderInformation(root, PreOrder, &writer)

		writer.Flush()

		if err != nil {
			t.Error(err)
		}
		reader := bitstream.CreateBitReader(writer.Bytes())

		treeFromHeader, err := createTreeFromHeader(&reader, size)

		if err != nil {
			t.Error(err)
		}

		got := make([]rune, 0)

		getInorderChars(treeFromHeader, &got)

		if !slices.Equal(expected, got) {
			t.Errorf("wrong tree created when reading the header")
		}

	})
}
