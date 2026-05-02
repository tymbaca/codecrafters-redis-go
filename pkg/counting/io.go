package counting

import "io"

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

type Writer struct {
	w io.Writer
	N int
}

func (w *Writer) Write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	w.N += n
	return n, err
}
