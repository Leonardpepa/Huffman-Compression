package bitstream

import (
	"testing"
	"unicode/utf8"
)

func btoi(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

func checkerr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
func TestReadBit(t *testing.T) {

	t.Run("read the bits 1 by 1", func(t *testing.T) {
		b := make([]byte, 0)

		b = append(b, byte(0b01101001))

		reader := CreateBitReader(b)

		bit1, err := reader.Read()
		checkerr(t, err)

		bit2, err := reader.Read()
		checkerr(t, err)

		bit3, err := reader.Read()
		checkerr(t, err)

		bit4, err := reader.Read()
		checkerr(t, err)

		bit5, err := reader.Read()
		checkerr(t, err)

		bit6, err := reader.Read()
		checkerr(t, err)

		bit7, err := reader.Read()
		checkerr(t, err)

		bit8, err := reader.Read()
		checkerr(t, err)

		if bit1 {
			t.Errorf("expected bit1: 0, got %d", btoi(bit1))
		}

		if !bit2 {
			t.Errorf("expected bit2: 1, got %d", btoi(bit2))
		}

		if !bit3 {
			t.Errorf("expected bit3: 1, got %d", btoi(bit3))
		}

		if bit4 {
			t.Errorf("expected bit4: 0, got %d", btoi(bit4))
		}

		if !bit5 {
			t.Errorf("expected bit5: 1, got %d", btoi(bit5))
		}

		if bit6 {
			t.Errorf("expected bit6: 0, got %d", btoi(bit6))
		}

		if bit7 {
			t.Errorf("expected bit7: 0, got %d", btoi(bit7))
		}

		if !bit8 {
			t.Errorf("expected bit8: 1, got %d", btoi(bit8))
		}

		bit9, err := reader.Read()

		if err == nil {
			t.Errorf("expected no more bits,  got %d", btoi(bit9))
		}

	})
}

func TestReadByte(t *testing.T) {
	t.Run("read the 1byte", func(t *testing.T) {
		b := make([]byte, 0)

		b = append(b, byte(0b11101001))

		reader := CreateBitReader(b)

		byte1, err := reader.ReadByte()

		checkerr(t, err)

		if byte1 != reader.Bytes()[0] {
			t.Errorf("expected 11101001, got %08b", byte1)
		}
	})
}

func TestReadCharacter(t *testing.T) {
	t.Run("read utf8 character", func(t *testing.T) {
		b := make([]byte, utf8.RuneLen(rune(256)))
		utf8.EncodeRune(b, rune(256))

		reader := CreateBitReader(b)

		expect, _ := utf8.DecodeRune(reader.Bytes())

		char, err := reader.ReadChar()

		checkerr(t, err)

		if char != expect {
			t.Errorf("expected %b, got %b", expect, char)
		}
	})
}

func TestReadBitAtPosition(t *testing.T) {
	t.Run("read the bit at position i", func(t *testing.T) {
		b := make([]byte, 0)

		b = append(b, byte(0b00100001))

		reader := CreateBitReader(b)

		bit, err := reader.getBitAt(reader.Bytes()[0], 0)

		checkerr(t, err)

		if bit {
			t.Errorf("expected 0, got %d", btoi(bit))
		}

		bit, err = reader.getBitAt(reader.Bytes()[0], 2)

		checkerr(t, err)

		if !bit {
			t.Errorf("expected 1, got %d", btoi(bit))
		}

		bit, err = reader.getBitAt(reader.Bytes()[0], 7)

		checkerr(t, err)

		if !bit {
			t.Errorf("expected 1, got %d", btoi(bit))
		}
	})
}

func TestHasNext(t *testing.T) {
	t.Run("should read (8 * total bytes) bits", func(t *testing.T) {
		b := make([]byte, 0)

		b = append(b, byte(0b00100000))

		reader := CreateBitReader(b)
		count := 0
		for reader.HasNext() {
			_, err := reader.Read()

			checkerr(t, err)
			count++
		}

		if count != reader.SizeOfBits() {
			t.Errorf("expeted count %d, got %d", reader.SizeOfBits(), count)
		}
	})
}
