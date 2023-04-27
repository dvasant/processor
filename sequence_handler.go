package processor

type SequenceHandler struct {
	input   chan *Envelope
	output  chan *Envelope
	counter int
	buffer  map[int]*Envelope
}

func NewSequenceHandler(input, output chan *Envelope) *SequenceHandler {
	return &SequenceHandler{
		input:   input,
		output:  output,
		counter: initialSqquenceValue,
		buffer:  make(map[int]*Envelope),
	}
}

func (this *SequenceHandler) Handle() {
	for envelope := range this.input {
		this.processEnvelope(envelope)
	}
}

func (this *SequenceHandler) processEnvelope(envelope *Envelope) {
	this.buffer[envelope.Sequence] = envelope
	this.sendBufferedEnelopersInOrder()
}

func (this *SequenceHandler) sendBufferedEnelopersInOrder() {
	for {
		next, found := this.buffer[this.counter]
		if !found {
			break
		}
		this.output <- next
		delete(this.buffer, this.counter)
		this.counter++
	}
}
