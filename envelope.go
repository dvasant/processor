package processor

const (
	initialSqquenceValue = 0
)

type (
	Envelope struct {
		EOF      bool
		Sequence int

		Input  AddressInput
		Output AddressOutput
	}

	AddressInput struct {
		Street1 string
		City    string
		State   string
		ZIPCode string
	}

	AddressOutput struct {
		Status        string
		DeliveryLine1 string
		LastLine      string
		City          string
		State         string
		ZIPCode       string
	}
)
