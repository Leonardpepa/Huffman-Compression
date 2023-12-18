package main

type BitWriter struct {
	data   []byte
	offset int
	count  int
}

func CreateBitWriter() BitWriter {
	return BitWriter{
		data: make([]byte, 0),
	}
}
