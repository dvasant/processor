package processor

import (
	"github.com/smartystreets/gunit"
	"net/http"
	"testing"
)

func TestVerifierFixture(t *testing.T) {
	gunit.Run(new(VerifierFixture), t)
}

type VerifierFixture struct {
	*gunit.Fixture
	client   *FakeHTTPClient
	verifier *SmartyVerifier
}

func (this *VerifierFixture) Setup() {
	this.client = &FakeHTTPClient{}
	this.verifier = NewSmartyVerifier(this.client)
}

func NewSmartyVerifier(client HTTPClient) *SmartyVerifier {
	return &SmartyVerifier{
		client: client,
	}
}

func (this *VerifierFixture) TestRequestComposedProperly() {
	input := AddressInput{
		Street1: "Street1",
		City:    "City",
		State:   "State",
		ZIPCode: "ZIPCode",
	}
	this.verifier.Verify(input)
	this.AssertEqual("GET", this.client.request.Method)
	this.AssertEqual("/street-address", this.client.request.URL.Path)
	this.AssertQueryStringValue("street", "Street1")
	this.AssertQueryStringValue("city", "City")
	this.AssertQueryStringValue("state", "State")
	this.AssertQueryStringValue("zipcode", "ZIPCode")
}

func (this *VerifierFixture) rawQuery() string {
	return this.client.request.URL.RawQuery
}

func (this *VerifierFixture) AssertQueryStringValue(key, expected string) {
	query := this.client.request.URL.Query()
	this.AssertEqual(expected, query.Get(key))
}

/////////////////////////////////////////////////////////

type FakeHTTPClient struct {
	request *http.Request
}

func (this *FakeHTTPClient) Do(request *http.Request) (*http.Response, error) {

	this.request = request
	return nil, nil
}
