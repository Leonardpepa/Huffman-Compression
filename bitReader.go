package main

import (
	"log"
)

type BitReader struct {
	data   []byte
	offset int
	count  int
	length int
}

func CreateBitReader(data []byte) BitReader {
	return BitReader{data: data, length: len(data)}
}

func (reader *BitReader) Bytes() []byte {
	return reader.data
}

func (reader *BitReader) Size() int {
	return reader.length
}

func (reader *BitReader) SizeOfBits() int {
	return reader.length * 8
}

func (reader *BitReader) Offset() int {
	return reader.offset
}

func (reader *BitReader) BitPosition() int {
	return reader.count
}

func (reader *BitReader) getBitAsBool(num byte, i int) bool {
	if i < 0 || i > 7 {
		log.Fatal("Error index must be in bounds [0-7]")
	}
	if num&(1<<i) == 0 {
		return false
	}

	return true
}

func (reader *BitReader) HasNext() bool {
	return reader.offset < reader.length
}

func (reader *BitReader) GetNextBit() bool {
	if reader.HasNext() == false {
		log.Fatal("No more bits to read. Abort")
	}

	bit := reader.getBitAsBool(reader.data[reader.offset], 7-reader.count)

	reader.count++
	if reader.count == 8 {
		reader.count = 0
		reader.offset++
	}

	return bit
}
