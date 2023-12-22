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

func (reader *Reader) getBitAt(num byte, i int) (bool, error) {
	if i < 0 || i > 7 {
		return false, fmt.Errorf("index out of bounds when reading bit")
	}
	return num&(1<<(7-i)) != 0, nil
}

func (reader *Reader) HasNext() bool {
	if len(reader.Bytes()) == 0 {
		return false
	}

	return reader.offset < reader.length
}

func (reader *Reader) Read() (bool, error) {

	if len(reader.Bytes()) == 0 {
		return false, fmt.Errorf("buffer is empty")
	}

	if reader.HasNext() == false {
		return false, fmt.Errorf("no more bits to read")
	}

	bit, err := reader.getBitAt(reader.data[reader.offset], reader.count)

	if err != nil {
		return false, err
	}
	reader.count++
	if reader.count == 8 {
		reader.count = 0
		reader.offset++
	}

	return bit, nil
}

// reads utf8 rune
func (reader *Reader) ReadChar() (rune, error) {
	if len(reader.Bytes()) == 0 {
		return 0, fmt.Errorf("buffer is empty")
	}

	byte1, err := reader.ReadByte()

	if err != nil {
		return 0, err
	}

	bit1, err := reader.getBitAt(byte1, 0)
	if err != nil {
		return 0, err
	}

	//0xxxxxxx
	if !bit1 {
		r, _ := utf8.DecodeRune([]byte{byte1})
		return r, nil
	}

	bit2, err := reader.getBitAt(byte1, 1)
	if err != nil {
		return 0, err
	}

	bit3, err := reader.getBitAt(byte1, 2)
	if err != nil {
		return 0, err
	}

	// 110xxxxx 10xxxxxx
	if bit1 && bit2 && !bit3 {
		byte2, err := reader.ReadByte()

		if err != nil {
			return 0, err
		}

		bit1, err = reader.getBitAt(byte2, 0)
		if err != nil {
			return 0, err
		}

		bit2, err = reader.getBitAt(byte2, 1)
		if err != nil {
			return 0, err
		}

		if !bit1 || bit2 {
			log.Printf("%08b, %08b", byte1, byte2)
			return 0, fmt.Errorf("error while decoding 110xxxxx 10xxxxxx utf8 rune")
		}

		r, _ := utf8.DecodeRune([]byte{byte1, byte2})
		return r, nil

	}

	bit4, err := reader.getBitAt(byte1, 3)
	if err != nil {
		return 0, err
	}

	// 1110xxxx 10xxxxxx 10xxxxxx
	if bit1 && bit2 && bit3 && !bit4 {
		byte2, err := reader.ReadByte()

		if err != nil {
			return 0, err
		}

		bit1, err = reader.getBitAt(byte2, 0)
		if err != nil {
			return 0, err
		}

		bit2, err = reader.getBitAt(byte2, 1)
		if err != nil {
			return 0, err
		}

		if !bit1 || bit2 {
			log.Printf("%08b, %08b, error in the second byte", byte1, byte2)
			return 0, fmt.Errorf("error while decoding 1110xxxx 10xxxxxx 10xxxxxx utf8 rune")
		}

		byte3, err := reader.ReadByte()

		if err != nil {
			return 0, err
		}

		bit1, err = reader.getBitAt(byte3, 0)
		if err != nil {
			return 0, err
		}

		bit2, err = reader.getBitAt(byte3, 1)
		if err != nil {
			return 0, err
		}

		if !bit1 || bit2 {
			log.Printf("%08b, %08b, %08b, error in the third byte", byte1, byte2, byte3)
			return 0, fmt.Errorf("error while decoding 1110xxxx 10xxxxxx 10xxxxxx utf8 rune")
		}

		r, _ := utf8.DecodeRune([]byte{byte1, byte2, byte3})
		return r, nil
	}

	// 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx
	bit5, err := reader.getBitAt(byte1, 4)
	if err != nil {
		return 0, err
	}

	if bit1 && bit2 && bit3 && bit4 && !bit5 {
		byte2, err := reader.ReadByte()

		if err != nil {
			return 0, err
		}

		bit1, err = reader.getBitAt(byte2, 0)
		if err != nil {
			return 0, err
		}

		bit2, err = reader.getBitAt(byte2, 1)
		if err != nil {
			return 0, err
		}

		if !bit1 || bit2 {
			log.Printf("%08b, %08b, error in the second byte", byte1, byte2)
			return 0, fmt.Errorf("error while decoding 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx utf8 rune")
		}

		byte3, err := reader.ReadByte()

		if err != nil {
			return 0, err
		}

		bit1, err = reader.getBitAt(byte3, 0)
		if err != nil {
			return 0, err
		}

		bit2, err = reader.getBitAt(byte3, 1)
		if err != nil {
			return 0, err
		}

		if !bit1 || bit2 {
			log.Printf("%08b, %08b, %08b, error in the fourth byte", byte1, byte2, byte3)
			return 0, fmt.Errorf("error while decoding 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx utf8 rune")
		}

		byte4, err := reader.ReadByte()

		if err != nil {
			return 0, err
		}

		bit1, err = reader.getBitAt(byte4, 0)
		if err != nil {
			return 0, err
		}

		bit2, err = reader.getBitAt(byte4, 1)
		if err != nil {
			return 0, err
		}

		if !bit1 || bit2 {
			log.Printf("%08b, %08b, %08b, %08b, error in the fourth byte", byte1, byte2, byte3, byte4)
			return 0, fmt.Errorf("error while decoding 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx utf8 rune")
		}

		r, _ := utf8.DecodeRune([]byte{byte1, byte2, byte3, byte4})
		return r, nil
	}

	return 0, fmt.Errorf("error while decoding utf8 rune")
}

func (reader *Reader) ReadByte() (byte, error) {
	if len(reader.Bytes()) == 0 {
		return 0, fmt.Errorf("buffer is empty")
	}

	writer := CreateBitWriter()
	for reader.HasNext() {
		bit, err := reader.Read()
		if err != nil {
			return 0, err
		}
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
