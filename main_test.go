package main

import (
	"huffmanCompression/huffman"
	"io"
	"log"
	"os"
	"slices"
	"testing"
)

func TestEncodeDecodeFile(t *testing.T) {
	t.Run("encode and decode file", func(t *testing.T) {

		inputFileName := "tests/hamlet.txt"
		outputFileName := "tests/decoded.txt"

		input, err := os.Open(inputFileName)
		if err != nil {
			t.Error(err)
		}

		defer func(input *os.File) {
			err := input.Close()
			if err != nil {
				log.Println(err)
			}
		}(input)

		err = huffman.Encode(input, inputFileName+".hf")
		if err != nil {
			t.Error(err)
		}

		encodedFile, err := os.Open(inputFileName + ".hf")
		if err != nil {
			t.Error(err)
		}

		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				log.Println(err)
			}
		}(inputFileName + ".hf")

		err = huffman.Decode(encodedFile, outputFileName)
		if err != nil {
			t.Error(err)
		}

		decoded, err := os.Open(outputFileName)
		if err != nil {
			t.Error(err)
		}

		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				log.Println(err)
			}
		}(outputFileName)

		_, err = input.Seek(0, io.SeekStart)
		if err != nil {
			t.Error(err)
		}

		expected, err := io.ReadAll(input)

		if err != nil {
			t.Error(err)
		}

		got, err := io.ReadAll(decoded)

		if err != nil {
			t.Error(err)
		}

		if !slices.Equal(expected, got) {
			t.Errorf("files dont match")
		}
	})
}
