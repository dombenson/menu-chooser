package courses

import (
	"database/sql"
	"encoding/json"
)

type Course interface {
	Name() string
	ID() int
	json.Marshaler
}

type course struct {
	id      int
	name    string
	eventid int
	order   int
}

const GetCoursesSQL = "SELECT `courseid`,`eventid`,`name`,`order` FROM `courses` WHERE `eventid` = ? ORDER BY `order`,`courseid` ASC"

func MakeCourse(rows *sql.Rows) (newCourse Course, err error) {
	retCourse := &course{}
	newCourse = retCourse
	err = rows.Scan(&retCourse.id, &retCourse.eventid, &retCourse.name, &retCourse.order)
	return
}

func (this *course) Name() string {
	return this.name
}
func (this *course) ID() int {
	return this.id
}

func (this *course) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		EventID  int    `json:"eventId"`
	}{
		ID:       this.id,
		Name:     this.name,
		EventID:  this.eventid,
	})
}
