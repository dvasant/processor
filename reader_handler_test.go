package processor

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"strconv"
	"testing"
)

func TestReaderHandlerFixture(t *testing.T) {
	gunit.Run(new(ReaderHandlerFixture), t)
}

type ReaderHandlerFixture struct {
	*gunit.Fixture
	buffer *ReadWriteSpyBuffer
	output chan *Envelope
	reader *ReaderHandler
}

func (this *ReaderHandlerFixture) Setup() {
	this.buffer = NewReadWriteSpyBuffer("")
	this.output = make(chan *Envelope, 10)
	this.reader = NewReaderHandler(this.buffer, this.output)

	const header = "Street1,City,State,ZIPCode"
	this.writeLine(header)
}

func (this *ReaderHandlerFixture) TestCSVRecordInEnvelope() {
	this.writeLine("A1,B1,C1,D1")
	this.reader.Handle()

	this.So(<-this.output, should.Resemble, buildEnvelope(initialSqquenceValue))
}

func (this *ReaderHandlerFixture) TestAllCSVRecordSentToOutput() {
	this.writeLine("A1,B1,C1,D1")
	this.writeLine("A2,B2,C2,D2")

	this.reader.Handle()

	this.assertRecordsSent()
	this.assertCleanUp()
}
func (this *ReaderHandlerFixture) writeLine(line string) {
	this.buffer.WriteString(line + "\n")
}

func (this *ReaderHandlerFixture) assertRecordsSent() {
	this.So(<-this.output, should.Resemble, buildEnvelope(initialSqquenceValue))
	this.So(<-this.output, should.Resemble, buildEnvelope(initialSqquenceValue+1))
}

func (this *ReaderHandlerFixture) assertCleanUp() {
	this.So(<-this.output, should.Resemble, &Envelope{Sequence: initialSqquenceValue + 2, EOF: true})
	this.So(<-this.output, should.BeNil)
	this.So(this.buffer.closed, should.Equal, 1)
}

func buildEnvelope(index int) *Envelope {
	suffix := strconv.Itoa(index + 1)
	return &Envelope{
		Input: AddressInput{
			Street1: "A" + suffix,
			City:    "B" + suffix,
			State:   "C" + suffix,
			ZIPCode: "D" + suffix,
		},
		Sequence: index,
	}
}

func (this *ReaderHandlerFixture) TestMalformedInputReturnsError() {
	malformedRecord := "A"
	this.writeLine(malformedRecord)
	err := this.reader.Handle()
	if this.So(err, should.NotBeNil) {
		this.So(err.Error(), should.Equal, "Malformed input")
	}
}
