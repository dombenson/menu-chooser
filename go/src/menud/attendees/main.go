package attendees

import (
	"crypto/hmac"
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
)

type Attendee interface {
	Name() string
	Email() string
	VerifyToken(string) error
	GetToken() string
}

type attendee struct {
	id      int
	name    string
	email   string
	key     string
	eventid int
}

const GetAttendeeSQL = "SELECT `attendeeid`,`name`,`email`,`key`,`eventid` FROM `attendees` WHERE `id` = ?"
const GetAttendeesByEventSQL = "SELECT `attendeeid` FROM `attendees` WHERE `eventid` = ?"

func MakeAttendee(rows *sql.Rows) (user Attendee, err error) {
	retUser := &attendee{}
	user = retUser
	err = rows.Scan(&retUser.id, &retUser.name, &retUser.email, &retUser.key, &retUser.eventid)
	return
}

func ParseToken(token string) (id int, authToken string) {
	byteArr := []byte(token)
	lenChar := byteArr[:1]
	intLen, err := strconv.Atoi(string(lenChar))
	if (err != nil) || (intLen < 1) || (intLen > 7) || (len(token) < (intLen + 1 + 12)) {
		return
	}
	idChars := byteArr[1:intLen]
	tmpId, err := strconv.Atoi(string(idChars))
	if (err != nil) || tmpId < 1 || tmpId > 9999999 {
		return
	}
	id = tmpId
	authToken = string(byteArr[(intLen + 1):])
	return
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
