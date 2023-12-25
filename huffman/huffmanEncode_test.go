package huffman

import (
	"bufio"
	"bytes"
	"huffmanCompression/bitstream"
	"testing"
)

func TestCreateBits(t *testing.T) {
	t.Run("should encode the content bit by bit", func(t *testing.T) {
		frequencies := map[rune]uint64{
			'M': 24,
			'C': 32,
			'K': 7,
			'D': 42,
			'E': 120,
			'L': 43,
			'U': 37,
			'Z': 2,
		}

		root := CreateHuffmanTreeFromFrequencies(frequencies)

		table := CalculateBitCodes(root)
		writer := bitstream.CreateBitWriter()

		text := "MKLULC"

		err := encodeBits(bufio.NewReader(bytes.NewReader([]byte(text))), &writer, table)

		if err != nil {
			t.Error(err)
		}

		writer.Flush()

		// expected bits
		expectedBits := ""
		for _, val := range text {
			expectedBits += table[val]
		}

		// 11111 111101 110 100 110 1110
		reader := bitstream.CreateBitReader(writer.Bytes())
		got := ""

		for reader.HasNext() {
			bit, err := reader.Read()

			if err != nil {
				t.Error(err)
			}
			if bit {
				got += string('1')
			} else {
				got += string('0')
			}
		}

		if expectedBits != got {
			t.Errorf("expected %s, got %s", expectedBits, got)
		}
	})
}
