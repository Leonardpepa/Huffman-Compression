package bitstream

import (
	"fmt"
	"unicode/utf8"
)

type Writer struct {
	data    []byte
	current byte
	count   uint
}

func CreateBitWriter() Writer {
	return Writer{
		data:    make([]byte, 0),
		current: 0,
		count:   0,
	}
}

func (writer *Writer) Bytes() []byte {
	return writer.data
}
func (writer *Writer) Size() int {
	return len(writer.data)
}

func (writer *Writer) BitPosition() uint {
	return writer.count
}

func (writer *Writer) WriteBitFromChar(bit rune) error {
	switch bit {
	case '1':
		writer.current = writer.current<<1 | 1
	case '0':
		writer.current = writer.current << 1
	default:
		return fmt.Errorf("Bit must be 0 or 1")
	}

	writer.count++

	if writer.count == 8 {
		writer.appendByte()
	}

	return nil
}

func (writer *Writer) WriteBitFromBool(bit bool) {
	switch bit {
	case true:
		writer.current = writer.current<<1 | 1
	case false:
		writer.current = writer.current << 1
	}

	writer.count++

	if writer.count == 8 {
		writer.appendByte()
	}
}
func (writer *Writer) appendByte() {
	writer.data = append(writer.data, writer.current)
	writer.current = byte(0)
	writer.count = 0
}

func (writer *Writer) WriteUtf8Rune(char rune) error {
	b := make([]byte, utf8.RuneLen(char))
	utf8.EncodeRune(b, char)

	bitReader := CreateBitReader(b)

	for bitReader.HasNext() {
		bit, err := bitReader.Read()
		if err != nil {
			return err
		}
		writer.WriteBitFromBool(bit)
	}
	return nil
}

func (writer *Writer) Flush() {
	if writer.HasRemainingBits() {
		writer.writeRemainingBitsWithPadding()
	} else {
		writer.current = byte(0)
		writer.count = 0
	}
}

func (writer *Writer) HasRemainingBits() bool {
	return writer.count != 0 && writer.count <= 7
}

func (writer *Writer) writeRemainingBitsWithPadding() {
	writer.data = append(writer.data, writer.current<<(8-writer.count))
	writer.current = byte(0)
	writer.count = 0
}
