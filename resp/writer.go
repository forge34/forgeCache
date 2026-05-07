package resp

import (
	"bufio"
	"io"
)

type Writer struct {
	Writer *bufio.Writer
}

func NewWriter(rd io.Writer) *Writer {
	return &Writer{Writer: bufio.NewWriter(rd)}
}

func (w *Writer) Write(v Value) error {
	p := v.Marshal()

	_, err := w.Writer.Write(p)
	if err != nil {
		return err
	}

	return w.Writer.Flush()
}
