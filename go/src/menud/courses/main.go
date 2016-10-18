package courses

import "database/sql"

type Course interface {
	Name() string
	ID() int
}

type course struct {
	id      int
	name    string
	eventid int
	order   int
}

const GetCoursesSQL = "SELECT `courseid`,`eventid`,`name`,`order` FROM `courses` WHERE `eventid` = ? ORDER BY order,courseid ASC"

func MakeCourse(rows *sql.Rows) (course Course, err error) {
	retCourse := &course{}
	course = retCourse
	err = rows.Scan(&retCourse.id, &retCourse.eventid, &retCourse.name, &retCourse.order)
	return
}

func (this *course) Name() string {
	return this.name
}
func (this *course) ID() string {
	return this.id
}
