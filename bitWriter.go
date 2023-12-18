package main

import (
	"log"
)

type BitWriter struct {
	data    []byte
	current byte
	count   uint
}

func CreateBitWriter() BitWriter {
	return BitWriter{
		data: make([]byte, 0),
	}
}

func (writer *BitWriter) Bytes() []byte {
	return writer.data
}
func (writer *BitWriter) Size() int {
	return len(writer.data)
}

func (writer *BitWriter) BitPosition() uint {
	return writer.count
}

func (writer *BitWriter) writeBitFromChar(bit rune) {
	switch bit {
	case '1':
		writer.current = writer.current<<1 | 1
	case '0':
		writer.current = writer.current << 1
	default:
		log.Fatal("Bit must be 0 or 1")
	}

	writer.count++

	if writer.count == 8 {
		writer.appendByte()
	}
}

func (writer *BitWriter) writeBitFromBool(bit bool) {
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
func (writer *BitWriter) appendByte() {
	writer.data = append(writer.data, writer.current)
	writer.current = byte(0)
	writer.count = 0
}

func (writer *BitWriter) writeRune(char rune) {
	bitReader := CreateBitReader([]byte(string(char)))

	for bitReader.HasNext() {
		bit := bitReader.Read()
		writer.writeBitFromBool(bit)
	}

	writer.WriteBytes()
}

func (writer *BitWriter) WriteBytes() {
	if writer.HasRemainingBits() {
		writer.WriteRemainingBitsWithPadding()
	}
}

func (writer *BitWriter) HasRemainingBits() bool {
	return writer.count%8 != 0
}

func (writer *BitWriter) WriteRemainingBitsWithPadding() {
	writer.data = append(writer.data, writer.current<<(8-writer.count))
	writer.current = byte(0)
	writer.count = 0
}
