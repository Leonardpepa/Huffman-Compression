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
func TestReadBit(t *testing.T) {

	b := make([]byte, 0)

	b = append(b, byte(0b01101001))

	reader := CreateBitReader(b)

	bit1 := reader.Read()
	bit2 := reader.Read()
	bit3 := reader.Read()
	bit4 := reader.Read()
	bit5 := reader.Read()
	bit6 := reader.Read()
	bit7 := reader.Read()
	bit8 := reader.Read()

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

	if err != nil {
		t.Error(err)
	}

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

	if err != nil {
		t.Error(err)
	}

	if char != expect {
		t.Errorf("expected %b, got %b", expect, char)
	}
}

func TestReadBitAtPosition(t *testing.T) {

	b := make([]byte, 0)

	b = append(b, byte(0b00100000))

	reader := CreateBitReader(b)

	bit := reader.getBitAt(reader.Bytes()[0], 2)

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
		reader.Read()
		count++
	}

	if count != 8 {
		t.Errorf("expeted count 8, got %d", count)
	}
}
