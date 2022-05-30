package swisspass_authenticator

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
	"net/http"
	"strconv"
	"time"
)

const CLIENT_ID_HEADER_KEY = "client_id"
const AUTHORIZATION_KEY = "Authorization"

type TokenVerifier interface {
	VerifyToken(clientId string, token string) (SwisspassUserInfo, error)
}

type SwisspassVerifier struct {
	endpoint    string
	restyclient *resty.Client
}

func New(endpoint string, timeoutInSeconds int) TokenVerifier {
	client := resty.New()
	client.SetTimeout(time.Duration(timeoutInSeconds) * time.Second)
	return SwisspassVerifier{endpoint, client}
}

func (sp SwisspassVerifier) VerifyToken(clientId string, token string) (SwisspassUserInfo, error) {
	var userInfo SwisspassUserInfo
	resp, err := sp.restyclient.R().
		SetHeader(CLIENT_ID_HEADER_KEY, clientId).
		SetHeader(AUTHORIZATION_KEY, token).
		Get(sp.endpoint + "/oev-oauth/oauth2-resource/" + clientId + "/userinfo")

	if err != nil {
		return userInfo, errors.Wrapf(err, "was not able to contact swisspass")
	}

	if resp.StatusCode() != http.StatusOK {
		return userInfo, errors.New("wrong or outdated token status code was: " + strconv.Itoa(resp.StatusCode()))
	}

	if err = json.Unmarshal(resp.Body(), &userInfo); err != nil {
		return userInfo, errors.Wrapf(err, "was not able to unmarshall json")
	}

	return userInfo, nil
}
