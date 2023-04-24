package processor

import (
	"encoding/csv"
	"io"
)

type ReaderHandler struct {
	reader *csv.Reader
	output chan *Envelope
	closer io.Closer
}

func NewReaderHandler(reader io.ReadCloser, output chan *Envelope) *ReaderHandler {
	return &ReaderHandler{
		reader: csv.NewReader(reader),
		closer: reader,
		output: output,
	}
}

func (this *ReaderHandler) Handle() {

}
