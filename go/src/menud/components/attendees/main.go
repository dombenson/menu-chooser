package attendees

import (
	"database/sql"
	"encoding/json"
	"strconv"
)

type Attendee interface {
	ID() int
	EventId() int
	Name() string
	Email() string
	VerifyToken(string,int) error
	GetToken() string
	GetLoginURL() string
	json.Marshaler
}

const GetAttendeeSQL = "SELECT `attendeeid`,`name`,`email`,`loginkey`,`eventid` FROM `attendees` WHERE `attendeeid` = ?"
const GetAttendeesByEventSQL = "SELECT `attendeeid` FROM `attendees` WHERE `eventid` = ?"

func MakeAttendee(rows *sql.Rows) (user Attendee, err error) {
	retUser := &attendee{}
	user = retUser
	err = rows.Scan(&retUser.id, &retUser.name, &retUser.email, &retUser.key, &retUser.eventid)
	return
}

func ParseToken(token string) (id int, authToken string, version int) {
	byteArr := []byte(token)
	verChar := string(byteArr[:1])
	offset := 0
	if verChar == "v" {
		verNumChar := byteArr[1:2]
		verNum, err := strconv.Atoi(string(verNumChar))
		if (err != nil) || (verNum < 1) || (verNum > 1) {
			// Only v1 currently supported
			return
		}
		version = verNum
		offset = 2
	}
	lenChar := byteArr[offset:offset+1]
	intLen, err := strconv.Atoi(string(lenChar))
	if (err != nil) || (intLen < 1) || (intLen > 7) || (len(token) < (intLen + 1 + 8)) {
		return
	}
	idChars := byteArr[offset+1 : offset+intLen+1]
	tmpId, err := strconv.Atoi(string(idChars))
	if (err != nil) || tmpId < 1 || tmpId > 9999999 {
		return
	}
	id = tmpId
	authToken = string(byteArr[(offset+intLen + 1):])
	return
}
