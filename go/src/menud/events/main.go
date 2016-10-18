package events

import (
	"database/sql"
	"time"
)

type Event interface {
	Name() string
	Location() string
	ID() int
	UserID() int
	Date() time.Time
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

func MakeEvent(rows *sql.Rows) (event Event, err error) {
	retEvent := &event{}
	event = retEvent
	err = rows.Scan(&retEvent.id, &retEvent.userid, &retEvent.name, &retEvent.location, &retEvent.date)
	return
}

func (this *event) Name() string {
	return this.name
}
func (this *event) Location() string {
	return this.location
}
func (this *event) ID() string {
	return this.id
}
func (this *event) UserID() string {
	return this.userid
}
func (this *event) Date() time.Time {
	return this.date
}
