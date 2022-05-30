package swisspass_authenticator

import (
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

const ISO_8601_DATE_FORMAT = "2006-02-01"

type SwisspassUserInfo struct {
	SwissPassId  string      `json:"SPIdPUID,omitempty"`
	Title        string      `json:"title,omitempty"`
	Salutation   string      `json:"salutation,omitempty"`
	FirstName    string      `json:"givenname,omitempty"`
	LastName     string      `json:"sn,omitempty"`
	ContactEmail string      `json:"contactEmail,omitempty"`
	AuthenEmail  string      `json:"authenEmail,omitempty"`
	Language     string      `json:"displayLanguage,omitempty"`
	Street       string      `json:"street,omitempty"`
	City         string      `json:"l,omitempty"`
	PostalCode   string      `json:"postalCode,omitempty"`
	CountryCode  string      `json:"c,omitempty"`
	Birthdate    *SimpleDate `json:"birthdate,omitempty"`
	Gender       string      `json:"gender,omitempty"`
	Tkid         string      `json:"tkid,omitempty"`
	CkmNumber    string      `json:"ckmNumber,omitempty"`
	CareOf       string      `json:"careOf,omitempty"`
	Postbox      string      `json:"postbox,omitempty"`
	Roles        []string    `json:"roles,omitempty"`
}

type SimpleDate struct {
	time.Time
}

func (t *SimpleDate) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)
	jsonTime, err := time.Parse(ISO_8601_DATE_FORMAT, string(strInput))
	if err != nil {
		log.Debugf("%s not parsable as date", string(strInput))
		return nil
	}
	*t = SimpleDate{jsonTime}
	return nil
}
