package processor

import (
	"fmt"
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
	this.sendEnvelopersInSequence(0, 1, 2, 3, 4)
	this.handler.Handle()
	this.So(this.sequenceOrder(), should.Resemble, []int{0, 1, 2, 3})
	if this.So(this.handler.buffer, should.BeEmpty) {
		fmt.Println(this.handler.buffer[6])
	}

}

func (this *SequenceHandlerFixture) TestEnvelopesReceivedOutofOrder_BufferedUtilsConfiguredBlock() {
	this.sendEnvelopersInSequence(4, 2, 0, 3, 1)

	this.handler.Handle()

	this.So(this.sequenceOrder(), should.Resemble, []int{0, 1, 2, 3})
	this.So(this.handler.buffer, should.BeEmpty)
}

func (this *SequenceHandlerFixture) sendEnvelopersInSequence(sequences ...int) {
	max := maxInt(sequences)
	for _, sequence := range sequences {
		this.input <- &Envelope{Sequence: sequence, EOF: max == sequence}
	}
}

func maxInt(ints []int) (max int) {
	for _, value := range ints {
		if value > max {
			max = value
		}
	}
	return max
}

func (this *SequenceHandlerFixture) sequenceOrder() (order []int) {
	for envelope := range this.output {
		order = append(order, envelope.Sequence)
	}
	return order
}
