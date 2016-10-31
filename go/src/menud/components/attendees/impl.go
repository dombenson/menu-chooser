package attendees

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"menud/config"
)

type attendee struct {
	id      int
	name    string
	email   string
	key     string
	eventid int
}

func (this *attendee) ID() int {
	return this.id
}
func (this *attendee) EventId() int {
	return this.eventid
}
func (this *attendee) Name() string {
	return this.name
}
func (this *attendee) Email() string {
	return this.email
}
func (this *attendee) VerifyToken(authToken string, version int) error {
	if authToken == this.getAuthToken(version) {
		return nil
	}
	return errors.New("Token invalid")
}
func (this *attendee) getAuthToken(version int) string {
	mac := hmac.New(sha512.New, []byte(this.key))
	mac.Write([]byte(fmt.Sprintf("%d:%s", this.id, this.email)))
	bytes := mac.Sum(nil)
	if version == 0 {
		return base64.StdEncoding.EncodeToString(bytes[12:20])
	}
	if version == 1 {
		return base64.URLEncoding.EncodeToString(bytes[12:21])
	}
	panic(errors.New(fmt.Sprintf("Auth token version %d not defined", version)))
}
func (this *attendee) GetToken() string {
	useVersion := 1
	strId := fmt.Sprintf("%d", this.id)
	idLen := len(strId)
	return fmt.Sprintf("v%d%d%d%s", useVersion, idLen, this.id, this.getAuthToken(useVersion))
}

func (this *attendee) GetLoginURL() string {
	return fmt.Sprintf("%s%s/api/login/%s", config.ExternalHost(), config.PathPrefix(), this.GetToken())
}

func (this *attendee) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		EventID int    `json:"eventId"`
	}{
		ID:      this.id,
		Name:    this.name,
		Email:   this.email,
		EventID: this.eventid,
	})
}
