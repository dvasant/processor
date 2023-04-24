package processor

import (
	"github.com/smartystreets/gunit"
	"testing"
)

func TestReaderHandlerFixture(t *testing.T) {
	gunit.Run(new(ReaderHandlerFixture), t)
}

type ReaderHandlerFixture struct {
	*gunit.Fixture
}

func (this *ReaderHandlerFixture) Setup() {
}

func (this *ReaderHandlerFixture) TestHeaderSkipped() {
	buffer := NewReadWriteSpyBuffer("Street1,City,State,ZIPCode")
	output := make(chan *Envelope, 10)
	reader := NewReaderHandler(buffer, output)

	reader.Handle()

	//this.So()
}
