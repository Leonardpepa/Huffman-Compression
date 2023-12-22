package bitstream

import (
	"testing"
	"unicode/utf8"
)

func btoi(b bool) int {
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
}

func TestReadByte(t *testing.T) {

	b := make([]byte, 0)

	b = append(b, byte(0b11101001))

	reader := CreateBitReader(b)

	byte1, err := reader.ReadByte()

	checkerr(t, err)

	if byte1 != reader.Bytes()[0] {
		t.Errorf("expected 11101001, got %08b", byte1)
	}
}

func TestReadCharacter(t *testing.T) {

	b := make([]byte, utf8.RuneLen(rune(256)))
	utf8.EncodeRune(b, rune(256))

	reader := CreateBitReader(b)

	expect, _ := utf8.DecodeRune(reader.Bytes())

	char, err := reader.ReadChar()

	checkerr(t, err)

	if char != expect {
		t.Errorf("expected %b, got %b", expect, char)
	}
}

func TestReadBitAtPosition(t *testing.T) {

	b := make([]byte, 0)

	b = append(b, byte(0b00100000))

	reader := CreateBitReader(b)

	bit, err := reader.getBitAt(reader.Bytes()[0], 2)

	checkerr(t, err)

	if !bit {
		t.Errorf("expected 1, got %d", btoi(bit))
	}
}

func TestHasNext(t *testing.T) {
	b := make([]byte, 0)

	b = append(b, byte(0b00100000))

	reader := CreateBitReader(b)
	count := 0
	for reader.HasNext() {
		_, err := reader.Read()

		checkerr(t, err)
		count++
	}

	if count != 8 {
		t.Errorf("expeted count 8, got %d", count)
	}
}
