package main

import (
	"bufio"
	"io"
)

var PseudoEOF = '\U0000FDEF'

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
	frequencies[PseudoEOF] = 0
}
