package processor

import (
	"github.com/smartystreets/gunit"
	"strings"
	"testing"
)

func TestHandlerFixture(t *testing.T) {
	gunit.Run(new(HandlerFixture), t)
}

type HandlerFixture struct {
	*gunit.Fixture

	input       chan *Envelope
	output      chan *Envelope
	application *FakeVerifyer
	handler     *VerifyHandler
}

func (this *HandlerFixture) Setup() {
	this.input = make(chan *Envelope, 10)
	this.output = make(chan *Envelope, 10)
	this.application = NewFakeVerifier()
	this.handler = NewVerifierHandler(this.input, this.output, this.application)
}

func (this *HandlerFixture) TestVerifierReceivesInput() {

	envelope := this.enqueueEnvelope("street")
	close(this.input)
	this.handler.Handle()

	this.AssertEqual("STREET", envelope.Output.DeliveryLine1)
	this.AssertEqual(envelope, <-this.output)
}
func (this *HandlerFixture) enqueueEnvelope(street1 string) *Envelope {
	envelope := &Envelope{
		Input: AddressInput{Street1: street1},
	}
	this.input <- envelope
	return envelope
}

func (this *HandlerFixture) TestInputQueueDrained() {
	envelope1 := this.enqueueEnvelope("41")
	envelope2 := this.enqueueEnvelope("43")
	envelope3 := this.enqueueEnvelope("44")
	close(this.input)

	this.handler.Handle()
	this.AssertEqual(envelope1, <-this.output)
	this.AssertEqual(envelope2, <-this.output)
	this.AssertEqual(envelope3, <-this.output)

}

// #####################################################

type FakeVerifyer struct {
	input  AddressInput
	output AddressOutput
}

func NewFakeVerifier() *FakeVerifyer {
	return &FakeVerifyer{}
}

func (this *FakeVerifyer) Verify(value AddressInput) AddressOutput {
	this.input = value
	return AddressOutput{DeliveryLine1: strings.ToUpper(value.Street1)}

}
