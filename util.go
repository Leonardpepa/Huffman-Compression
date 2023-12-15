package main

import (
	"bufio"
	"io"
	"math"
)

var PseudoEOF = rune(math.MaxInt32)

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

	if frequencies[PseudoEOF] == 0 {
		frequencies[PseudoEOF] = 0
	}

	return frequencies, nil
}
