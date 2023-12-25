# Huffman Compression written in go

## Purpose
This project is a solution for [Write Your Own Compression Tool](https://codingchallenges.fyi/challenges/challenge-huffman)
build for my personal educational purposes

# Description
This is an implementation of the huffman compression algorithm. Read how it works here [Huffman Coding Trees](https://opendsa-server.cs.vt.edu/ODSA/Books/CS3/html/Huffman.html)

# Features
* Compress file
* Decompress fie

# Implementation Details
## header - Huffman Topology
In order to decode the file we need the original tree that was used to encode the file. This implementation
writes the tree structure in the file header using preorder traversal, writing 0 for each internal node
and 1 followed by the character for each leaf node. The size of the header is stored first followed by the actual information needed to created the tree

## Pseudo EOF
Because the operate in byte chunks and this algorithm writes the encoded content bit by bit
we need a way to know when the file contents end, inorder to avoid reading any junk bits added to fill the last byte.
This implementation uses the technique of the Pseudo EOF. The fixed-length code of the Pseudo EOF is stored at the end of the file. 
The program stops reading when it encounters the Pseudo EOF. The Pseudo EOF is a character that doesn't occur in the original file.

## Bit Stream
The bit strings of the fixed-length codes for the characters need to be stored bit by bit in the file.
Because go doesn't provide such functionality, bitstream package was implemented from scratch inorder to be able to read and write bits, bytes and runes bit by bit.

# Usage
```terminal
Usage: ./huffmanCompression.exe [OPTION]... [FILE]
[OPTIONS]:
-c FILE ~file to compress
-d FILE ~file to decompress
[FILE]:
the file specified is the file to compress or decompress,
when compressed the file will be saved with a .hf extension,
when decompressed the file will be saved as a copy of the original file
```

# Example 
## Compress
```terminal
.\huffmanCompression.exe -c .\tests\gutenberg.txt
2023/12/25 02:52:09 Opening file...  .\tests\gutenberg.txt
2023/12/25 02:52:09 Creating huffman tree...
2023/12/25 02:52:09 Calculating frequencies...
2023/12/25 02:52:09 Create priority queue...
2023/12/25 02:52:09 Build the tree...
2023/12/25 02:52:09 Calculating variable length codes...
2023/12/25 02:52:09 Encoding the tree in the header...
2023/12/25 02:52:09 Encoding the data bit by bit...
2023/12/25 02:52:09 Encoding finished. File: gutenberg.txt.hf,
Original: (3369045 bytes) compressed: (1919959 bytes)
```
## Decompress
```terminal
.\huffmanCompression.exe -d .\gutenberg.txt.hf
2023/12/25 02:53:52 Opening file...  .\gutenberg.txt.hf
2023/12/25 02:53:52 Reading encoded data... 
2023/12/25 02:53:52 Constructing the huffman tree from the header... 
2023/12/25 02:53:52 Decoding the data... 
2023/12/25 02:53:52 Decoding finished. File: gutenberg.txt, Bytes: 3369045
```

# How to run
1. Clone the repo ```git clone https://github.com/Leonardpepa/Huffman-Compression```
2. Build ```go build```
3. run on windows```huffmanCompression.exe [OPTION] [FILE]```
4. run on linux ```.\huffmanCompression [OPTION] [FILE]```
5. run all tests ```go test ./...```
