package main

import (
	"bufio"
	"io"
)

var PseudoEOF = rune(256)

func CalculateFrequencies(reader *bufio.Reader) (map[rune]uint64, error) {

	frequencies := make(map[rune]uint64)
	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		frequencies[char]++
	}
	return frequencies, nil
}

func AddPseudoEOF(frequencies map[rune]uint64) {
	for _, ok := frequencies[PseudoEOF]; ok; {
		PseudoEOF++
	}
	frequencies[PseudoEOF] = 0

}
