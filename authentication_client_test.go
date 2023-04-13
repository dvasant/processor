package processor

import (
	"errors"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticationClient(t *testing.T) {
	gunit.Run(new(AuthenticationClientFixture), t)
}

type AuthenticationClientFixture struct {
	*gunit.Fixture
	inner  *FakeHTTPClient
	client *AuthenticationClient
}

func (this *AuthenticationClientFixture) Setup() {
	this.inner = &FakeHTTPClient{}
	this.client = NewAuthenticationClient(this.inner, "https", "HOSTNAME", "authid", "authtoken")
}

func (this *AuthenticationClientFixture) TestProvidedInformationAddedBeforeRequestIsSent() {
	request := httptest.NewRequest("GET", "/path?existingKey=existingValue", nil)

	_, _ = this.client.Do(request)
	this.assertRequestConnectionInformation()
	this.assertQueryStringIncludesAuthentication()
}

func (this *AuthenticationClientFixture) TestResponseAndErrorFromInnerClientReturned() {
	this.inner.response = &http.Response{StatusCode: http.StatusTeapot}
	this.inner.err = errors.New("HTTP Error")
	request := httptest.NewRequest("GET", "/path", nil)
	response, err := this.client.Do(request)
	this.So(response.StatusCode, should.Equal, http.StatusTeapot)
	this.So(err.Error(), should.Equal, "HTTP Error")
}

func (this *AuthenticationClientFixture) assertQueryStringValue(key string, expectedValue string) {
	this.So(this.inner.request.URL.Query().Get(key), should.Equal, expectedValue)

}

func (this *AuthenticationClientFixture) assertQueryStringIncludesAuthentication() {
	this.assertQueryStringValue("auth-id", "authid")
	this.assertQueryStringValue("auth-token", "authtoken")
	this.assertQueryStringValue("existingKey", "existingValue")
}

func (this *AuthenticationClientFixture) assertRequestConnectionInformation() {

	this.So(this.inner.request.URL.Scheme, should.Equal, "https")
	this.So(this.inner.request.Host, should.Equal, "HOSTNAME")
	this.So(this.inner.request.URL.Host, should.Equal, "HOSTNAME")

}
