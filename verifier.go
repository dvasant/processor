package processor

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type SmartyVerifier struct {
	client HTTPClient
}

func (this *SmartyVerifier) Verify(input AddressInput) AddressOutput {

	response, _ := this.client.Do(this.buildRequest(input))
	candidates := this.decodeResponse(response)
	return this.prepareAddressOutput(candidates)
}

func (this *SmartyVerifier) buildRequest(input AddressInput) *http.Request {
	query := make(url.Values)
	query.Set("street", input.Street1)
	query.Set("city", input.City)
	query.Set("state", input.State)
	query.Set("zipcode", input.ZIPCode)
	request, _ := http.NewRequest("GET", "/street-address?"+query.Encode(), nil)
	return request
}

func (this *SmartyVerifier) decodeResponse(response *http.Response) (output []Candidate) {
	if response != nil {
		defer response.Body.Close()
		json.NewDecoder(response.Body).Decode(&output)
	}
	return output
}

func (this *SmartyVerifier) prepareAddressOutput(candidates []Candidate) AddressOutput {
	if len(candidates) == 0 {
		return AddressOutput{Status: "Invalid API Response"}
	}

	candidate := candidates[0]
	return AddressOutput{
		Status:        computeStatus(candidate),
		DeliveryLine1: candidate.DeliveryLine1,
		LastLine:      candidate.LastLine,
		City:          candidate.Components.City,
		State:         candidate.Components.State,
		ZIPCode:       candidate.Components.ZIPCode,
	}
}

func computeStatus(candidate Candidate) string {
	analysis := candidate.Analysis
	if !isDeliverable(analysis.Match) {
		return "Invalid"
	} else if analysis.Vacant == "Y" {
		return "Vacant"
	} else if analysis.Active != "Y" {
		return "Inactive"
	} else {
		return "Deliverable"
	}
}

func isDeliverable(value string) bool {
	return value == "Y" || value == "D" || value == "S"
}
