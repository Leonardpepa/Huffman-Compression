package bitstream

import (
	"testing"
	"unicode/utf8"
)

func TestWritingBits(t *testing.T) {
	t.Run("write bit from booleans", func(t *testing.T) {
		writer := CreateBitWriter()

		// write 01111111

		writer.WriteBitFromBool(false)
		writer.WriteBitFromBool(true)
		writer.WriteBitFromBool(true)
		writer.WriteBitFromBool(true)
		writer.WriteBitFromBool(true)
		writer.WriteBitFromBool(true)
		writer.WriteBitFromBool(true)
		writer.WriteBitFromBool(true)

		// here it doesn't do much because we already have written a full byte
		// here can be omitted
		writer.Flush()

		if len(writer.Bytes()) != 1 {
			t.Errorf("more or less than 1 byte have beem written")
		}

		if writer.Bytes()[0] != byte(0b01111111) {
			t.Errorf("expected %08b, got %08b", byte(0b01111111), writer.Bytes()[0])
		}

	})

	t.Run("write bits from chars '0', '1'", func(t *testing.T) {
		writer := CreateBitWriter()

		// write 11110111
		err := writer.WriteBitFromChar('1')
		checkerr(t, err)
		err = writer.WriteBitFromChar('1')
		checkerr(t, err)
		err = writer.WriteBitFromChar('1')
		checkerr(t, err)
		err = writer.WriteBitFromChar('1')
		checkerr(t, err)
		err = writer.WriteBitFromChar('0')
		checkerr(t, err)
		err = writer.WriteBitFromChar('1')
		checkerr(t, err)
		err = writer.WriteBitFromChar('1')
		checkerr(t, err)
		err = writer.WriteBitFromChar('1')
		checkerr(t, err)

		// here it doesn't do much because we already have written a full byte
		// here can be omitted
		writer.Flush()

		if len(writer.Bytes()) != 1 {
			t.Errorf("more or less than 1 byte have beem written")
		}

		if writer.Bytes()[0] != byte(0b11110111) {
			t.Errorf("expected %08b, got %08b", byte(0b11110111), writer.Bytes()[0])
		}

	})

}

func TestWriteUtf8RuneCharacter(t *testing.T) {
	t.Run("write a utf8 character", func(t *testing.T) {
		writer := CreateBitWriter()

		err := writer.WriteUtf8Rune('#')
		checkerr(t, err)

		char, _ := utf8.DecodeRune(writer.Bytes())

		if char != '#' {
			t.Errorf("expected #, got %c", char)
		}
	})
}

func TestWriter_Flush(t *testing.T) {
	t.Run("write 2 bits and pad the remaining with zeros", func(t *testing.T) {
		writer := CreateBitWriter()

		writer.WriteBitFromBool(false)
		writer.WriteBitFromBool(true)

		// writes remaining bits to complete a byte
		// fills the byte with trailing zeros
		// call only when you are finished writing individual bits
		writer.Flush()

		if len(writer.Bytes()) != 1 {
			t.Errorf("more or less than 1 byte have beem written")
		}

		if writer.Bytes()[0] != byte(0b01000000) {
			t.Errorf("expected %08b, got %08b", byte(0b01000000), writer.Bytes()[0])
		}

	})
}

func TestWriteCharAlongSideBits(t *testing.T) {
	t.Run("write 3 bits followed by a character (like huffman coding)", func(t *testing.T) {
		writer := CreateBitWriter()

		writer.WriteBitFromBool(true)
		writer.WriteBitFromBool(false)
		writer.WriteBitFromBool(true)

		err := writer.WriteUtf8Rune('^')
		checkerr(t, err)

		// be sure to write the remaining bits that don't complete a byte
		// [1010101111000000]
		// the first three bits we wrote individually
		// the next 8 are the binary code for the ^ character
		// the remaining are the padding added by calling the flush so we can complete the second byte
		writer.Flush()

		bitReader := CreateBitReader(writer.Bytes())

		bit1, err := bitReader.Read()
		checkerr(t, err)

		bit2, err := bitReader.Read()
		checkerr(t, err)

		bit3, err := bitReader.Read()
		checkerr(t, err)

		char, err := bitReader.ReadChar()
		checkerr(t, err)

		if !bit1 {
			t.Errorf("expected 1, got %d", btoi(bit1))
		}
		if bit2 {
			t.Errorf("expected 0, got %d", btoi(bit2))
		}
		if !bit3 {
			t.Errorf("expected 1, got %d", btoi(bit3))
		}
		if char != '^' {
			t.Errorf("expected ^, got %c", char)
		}

		if !bitReader.HasNext() || 8-bitReader.BitPosition() != 5 {
			t.Errorf("there should be  bits left to be read")
		}

	})
}
