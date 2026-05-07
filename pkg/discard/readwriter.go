package discard

import "io"

var Discard = readWriter{}

type readWriter struct{}

func (re readWriter) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (re readWriter) Write(p []byte) (n int, err error) {
	return io.Discard.Write(p)
}

func (re readWriter) Close() error {
	return nil
}
