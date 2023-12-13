package main

import (
	"bufio"
	"io"
)

func calculateFrequencies(reader *bufio.Reader) (map[rune]uint64, error) {
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
