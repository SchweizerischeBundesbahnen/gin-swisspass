package swisspass_authenticator

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

const testToken = "{\"birthdate\":\"1970-10-11\",\"c\":\"CH\",\"gender\":\"FEMALE\",\"tkid\":\"d6d20211-8b6d-4e2c-bc3b-bcf1a669effb\",\"contactEmail\":\"tokentest@grr.la\",\"postalCode\":\"1234\",\"displayLanguage\":\"de\",\"l\":\"asdas\",\"SPIdPUID\":\"fbfe3b6d-8485-4b92-aeb5-50a13baf5970\",\"authenEmail\":\"tokentest@grr.la\",\"givenname\":\"asd\",\"street\":\"ads\",\"lastDataUpdate\":\"2019-06-26 13:26:07\",\"salutation\":\"FRAU\",\"sn\":\"asd\"}"
const testSwisspassUrl = "https://www-test.swisspass.ch"
const GivenClientId = "oauth_tester_test"
const GivenToken = "Bearer fbfe3b6d-8485-4b92-aeb5-50a13baf5970.YL4qY1f0SlTASHCR_J4R0CzSsNIxeIAZ3ykafQPe"

func Test_VerifyTokenVerify(t *testing.T) {
	sp := New(testSwisspassUrl, 1).(SwisspassVerifier)
	httpmock.ActivateNonDefault(sp.restyclient.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testSwisspassUrl+"/oev-oauth/oauth2-resource/oauth_tester_test/userinfo", httpmock.NewStringResponder(200, testToken))

	userInfo, err := sp.VerifyToken(GivenClientId, GivenToken)

	assert.Nil(t, err)
	assert.NotNil(t, userInfo)
	expected, _ := time.Parse("2006-02-01", "1970-10-11")
	assert.True(t, userInfo.Birthdate.Time.Equal(expected))
	assert.Equal(t, "CH", userInfo.CountryCode)
	assert.Equal(t, "FEMALE", userInfo.Gender)
	assert.Equal(t, "d6d20211-8b6d-4e2c-bc3b-bcf1a669effb", userInfo.Tkid)
	assert.Equal(t, "tokentest@grr.la", userInfo.ContactEmail)
	assert.Equal(t, "1234", userInfo.PostalCode)
	assert.Equal(t, "de", userInfo.Language)
	assert.Equal(t, "asdas", userInfo.City)
	assert.Equal(t, "fbfe3b6d-8485-4b92-aeb5-50a13baf5970", userInfo.SwissPassId)
	assert.Equal(t, "tokentest@grr.la", userInfo.AuthenEmail)
	assert.Equal(t, "asd", userInfo.FirstName)
	assert.Equal(t, "asd", userInfo.LastName)
	assert.Equal(t, "FRAU", userInfo.Salutation)
	assert.Equal(t, "ads", userInfo.Street)
}

func Test_OutdatedToken(t *testing.T) {
	sp := New(testSwisspassUrl, 1).(SwisspassVerifier)
	httpmock.ActivateNonDefault(sp.restyclient.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testSwisspassUrl+"/oev-oauth/oauth2-resource/oauth_tester_test/userinfo", httpmock.NewStringResponder(400, testToken))

	_, err := sp.VerifyToken(GivenClientId, GivenToken)

	assert.NotNil(t, err)
	assert.Equal(t, "wrong or outdated token status code was: 400", err.Error())
}

func TestGinSwisspass_Headers_To_Swisspass_Are_Correct(t *testing.T) {
	sp := New(testSwisspassUrl, 1).(SwisspassVerifier)
	httpmock.ActivateNonDefault(sp.restyclient.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testSwisspassUrl+"/oev-oauth/oauth2-resource/oauth_tester_test/userinfo", CreateHeadersCheckResponder(t))
	_, _ = sp.VerifyToken(GivenToken, GivenToken)
}

func CreateHeadersCheckResponder(t *testing.T) httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		res := httpmock.NewStringResponse(200, testToken)
		res.Request = req

		assert.Equal(t, GivenToken, req.Header.Get(CLIENT_ID_HEADER_KEY))
		assert.Equal(t, GivenToken, req.Header.Get(AUTHORIZATION_KEY))

		return res, nil
	}
}
