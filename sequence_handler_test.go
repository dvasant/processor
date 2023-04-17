package processor

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"testing"
)

func TestSequenceHandler(t *testing.T) {
	gunit.Run(new(SequenceHandlerFixture), t)
}

type SequenceHandlerFixture struct {
	*gunit.Fixture
	input   chan *Envelope
	output  chan *Envelope
	handler *SequenceHandler
}

func (this *SequenceHandlerFixture) Setup() {
	this.input = make(chan *Envelope, 10)
	this.output = make(chan *Envelope, 10)
	this.handler = NewSequenceHandler(this.input, this.output)
}

func (this *SequenceHandlerFixture) TestExpectedEnvelopeSentToOutput() {
	this.sendEnvelopersInsequence(0, 1, 2, 3, 4)
	this.handler.Handle()
	this.So(this.sequenceOreder(), should.Resemble, []int{0, 1, 2, 3, 4})
}

func (this *SequenceHandlerFixture) TestEnvelopesReceivedOutofOrder_BufferedUtilsConfiguredBlock() {
	this.sendEnvelopersInsequence(4, 2, 0, 3, 1)

	this.handler.Handle()

	this.So(this.sequenceOreder(), should.Resemble, []int{0, 1, 2, 3, 4})
	this.So(this.handler.buffer, should.BeEmpty)
}

func (this *SequenceHandlerFixture) sendEnvelopersInsequence(sequences ...int) {
	for _, sequence := range sequences {
		this.input <- &Envelope{Sequence: sequence}
	}
	close(this.input)
}

func (this *SequenceHandlerFixture) sequenceOreder() (order []int) {

	close(this.output)
	for envelope := range this.output {
		order = append(order, envelope.Sequence)
	}
	return order
}
