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

func (reader *Reader) getBitAt(num byte, i int) bool {
	if i < 0 || i > 7 {
		log.Fatal("Error index must be in bounds [0-7]")
	}
	return num&(1<<(7-i)) != 0
}

func (reader *Reader) HasNext() bool {
	return reader.offset < reader.length
}

func (reader *Reader) Read() bool {
	if reader.HasNext() == false {
		log.Fatal("No more bits to read. Abort")
	}

	bit := reader.getBitAt(reader.data[reader.offset], reader.count)

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
		return 0, err
	}

	bit1 := reader.getBitAt(byte1, 0)
	//0xxxxxxx
	if bit1 == false {
		r, _ := utf8.DecodeRune([]byte{byte1})
		return r, nil
	}

	bit2 := reader.getBitAt(byte1, 1)
	bit3 := reader.getBitAt(byte1, 2)

	// 110xxxxx 10xxxxxx
	if bit1 && bit2 && bit3 == false {
		byte2, err := reader.ReadByte()

		if err != nil {
			return 0, err
		}

		bit1 = reader.getBitAt(byte2, 0)
		bit2 = reader.getBitAt(byte2, 1)

		if bit1 == false || bit2 == true {
			log.Printf("%08b, %08b", byte1, byte2)
			return 0, fmt.Errorf("error while decoding 110xxxxx 10xxxxxx utf8 rune")
		}

		r, _ := utf8.DecodeRune([]byte{byte1, byte2})
		return r, nil

	}

	bit4 := reader.getBitAt(byte1, 3)
	// 1110xxxx 10xxxxxx 10xxxxxx
	if bit1 && bit2 && bit3 && bit4 == false {
		byte2, err := reader.ReadByte()

		if err != nil {
			return 0, err
		}

		bit1 = reader.getBitAt(byte2, 0)
		bit2 = reader.getBitAt(byte2, 1)

		if bit1 == false || bit2 == true {
			log.Printf("%08b, %08b, error in the second byte", byte1, byte2)
			return 0, fmt.Errorf("error while decoding 1110xxxx 10xxxxxx 10xxxxxx utf8 rune")
		}

		byte3, err := reader.ReadByte()

		if err != nil {
			return 0, err
		}

		bit1 = reader.getBitAt(byte3, 0)
		bit2 = reader.getBitAt(byte3, 1)

		if bit1 == false || bit2 == true {
			log.Printf("%08b, %08b, %08b, error in the third byte", byte1, byte2, byte3)
			return 0, fmt.Errorf("error while decoding 1110xxxx 10xxxxxx 10xxxxxx utf8 rune")
		}

		r, _ := utf8.DecodeRune([]byte{byte1, byte2, byte3})
		return r, nil
	}

	// 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx
	bit5 := reader.getBitAt(byte1, 4)

	if bit1 && bit2 && bit3 && bit4 && bit5 == false {
		byte2, err := reader.ReadByte()

		if err != nil {
			return 0, err
		}

		bit1 = reader.getBitAt(byte2, 0)
		bit2 = reader.getBitAt(byte2, 1)

		if bit1 == false || bit2 == true {
			log.Printf("%08b, %08b, error in the second byte", byte1, byte2)
			return 0, fmt.Errorf("error while decoding 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx utf8 rune")
		}

		byte3, err := reader.ReadByte()

		if err != nil {
			return 0, err
		}

		bit1 = reader.getBitAt(byte3, 0)
		bit2 = reader.getBitAt(byte3, 1)

		if bit1 == false || bit2 == true {
			log.Printf("%08b, %08b, %08b, error in the fourth byte", byte1, byte2, byte3)
			return 0, fmt.Errorf("error while decoding 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx utf8 rune")
		}

		byte4, err := reader.ReadByte()

		if err != nil {
			return 0, err
		}

		bit1 = reader.getBitAt(byte4, 0)
		bit2 = reader.getBitAt(byte4, 1)

		if bit1 == false || bit2 == true {
			log.Printf("%08b, %08b, %08b, %08b, error in the fourth byte", byte1, byte2, byte3, byte4)
			return 0, fmt.Errorf("error while decoding 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx utf8 rune")
		}

		r, _ := utf8.DecodeRune([]byte{byte1, byte2, byte3, byte4})
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
