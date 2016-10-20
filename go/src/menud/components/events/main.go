package events

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Event interface {
	Name() string
	Location() string
	ID() int
	UserID() int
	Date() time.Time
	json.Marshaler
}

type event struct {
	id       int
	userid   int
	name     string
	location string
	date     time.Time
}

const GetEventsForUserSQL = "SELECT `eventid` FROM `events` WHERE `userid` = ? ORDER BY date ASC"
const GetEventSQL = "SELECT `eventid`,`userid`,`name`,`location`,`date` FROM `events` WHERE `eventid` = ?"

func MakeEvent(rows *sql.Rows) (newEvent Event, err error) {
	retEvent := &event{}
	newEvent = retEvent
	err = rows.Scan(&retEvent.id, &retEvent.userid, &retEvent.name, &retEvent.location, &retEvent.date)
	return
}

func (this *event) Name() string {
	return this.name
}
func (this *event) Location() string {
	return this.location
}
func (this *event) ID() int {
	return this.id
}
func (this *event) UserID() int {
	return this.userid
}
func (this *event) Date() time.Time {
	return this.date
}

func (this *event) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Location string `json:"location"`
		Date     string `json:"date"`
	}{
		ID:       this.id,
		Name:     this.name,
		Location: this.location,
		Date:     this.date.Format(time.RFC3339),
	})
}
