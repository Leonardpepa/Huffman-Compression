package huffman

import (
	"huffmanCompression/bitstream"
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
