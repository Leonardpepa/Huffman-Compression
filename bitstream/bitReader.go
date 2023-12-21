package bitstream

import (
	"fmt"
	"log"
	"unicode/utf8"
)

type Reader struct {
	data   []byte
	offset int
	count  int
	length int
}

func CreateBitReader(data []byte) Reader {
	return Reader{data: data, length: len(data)}
}

func (reader *Reader) Bytes() []byte {
	return reader.data
}

func (reader *Reader) Size() int {
	return reader.length
}

func (reader *Reader) SizeOfBits() int {
	return reader.length * 8
}

func (reader *Reader) Offset() int {
	return reader.offset
}

func (reader *Reader) BitPosition() int {
	return reader.count
}

func (reader *Reader) getBitAsBool(num byte, i int) bool {
	if i < 0 || i > 7 {
		log.Fatal("Error index must be in bounds [0-7]")
	}
	if num&(1<<i) == 0 {
		return false
	}

	return true
}

func (reader *Reader) HasNext() bool {
	return reader.offset < reader.length
}

func (reader *Reader) Read() bool {
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

func (reader *Reader) ReadChar() (rune, error) {

	byte1, err := reader.ReadByte()

	if err != nil {

	}
	//0xxxxxxx
	if reader.getBitAsBool(byte1, 7) == false {
		r, _ := utf8.DecodeRune([]byte{byte1})
		return r, nil
	}

	// 110xxxxx 10xxxxxx
	if reader.getBitAsBool(byte1, 7) &&
		reader.getBitAsBool(byte1, 6) &&
		reader.getBitAsBool(byte1, 5) == false {
		byte2, err := reader.ReadByte()

		if err != nil {
			return 0, err
		}

		r, _ := utf8.DecodeRune([]byte{byte1, byte2})
		return r, nil

	}

	// 1110xxxx 10xxxxxx 10xxxxxx
	if reader.getBitAsBool(byte1, 7) &&
		reader.getBitAsBool(byte1, 6) &&
		reader.getBitAsBool(byte1, 5) &&
		reader.getBitAsBool(byte1, 4) == false {
		byte2, err := reader.ReadByte()

		if err != nil {
			return 0, err
		}
		byte3, err := reader.ReadByte()

		if err != nil {
			return 0, err
		}
		r, _ := utf8.DecodeRune([]byte{byte1, byte2, byte3})
		return r, nil

	}

	return 0, fmt.Errorf("error while decoding utf8 rune")
}

func (reader *Reader) ReadByte() (byte, error) {
	writer := CreateBitWriter()
	for reader.HasNext() {
		bit := reader.Read()
		writer.WriteBitFromBool(bit)

		if writer.Size() == 1 {
			break
		}
	}

	// it may stop before it reads 8bits
	if writer.Size() == 0 {
		return 0, fmt.Errorf("There are not enough bits to read a byte. ")
	}

	return writer.Bytes()[0], nil
}