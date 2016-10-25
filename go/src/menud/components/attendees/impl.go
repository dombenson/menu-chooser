package attendees

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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
func (this *attendee) VerifyToken(authToken string) error {
	if authToken == this.getAuthToken() {
		return nil
	}
	return errors.New("Token invalid")
}
func (this *attendee) getAuthToken() string {
	mac := hmac.New(sha512.New, []byte(this.key))
	mac.Write([]byte(fmt.Sprintf("%d:%s", this.id, this.email)))
	bytes := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(bytes[12:20])
}
func (this *attendee) GetToken() string {
	strId := fmt.Sprintf("%d", this.id)
	idLen := len(strId)
	return fmt.Sprintf("%d%d%s", idLen, this.id, this.getAuthToken())
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
