package db

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"menud/attendees"
	"menud/config"
	"menud/courses"
	"menud/events"
	"menud/options"
	"menud/users"
)

type Connection interface {
	GetUser(int) (users.User, error)
	GetUserByEmailPassword(email, password string) (users.User, error)
}

type connection struct {
	baseConn              *sql.DB
	getUserStmt           *sql.Stmt
	getUserEmailStmt      *sql.Stmt
	getAttendeeStmt       *sql.Stmt
	getEventAttendeesStmt *sql.Stmt
	getUserEventsStmt     *sql.Stmt
	getEventStmt          *sql.Stmt
	getCoursesStmt        *sql.Stmt
	getOptionsStmt        *sql.Stmt
	getSelectionStmt      *sql.Stmt
}

func GetConnection() Connection {
	obj := &connection{}
	obj.connect()
	return obj
}

func (this *connection) connect() {
	var err error
	this.baseConn, err = sql.Open("mysql", config.DBConnString())
	if err != nil {
		panic("Unable to connect to database")
	}
	this.getUserStmt, err = this.baseConn.Prepare(users.GetUserSQL)
	this.getUserEmailStmt, err = this.baseConn.Prepare(users.GetUserByEmailSQL)
	this.getAttendeeStmt, err = this.baseConn.Prepare(attendees.GetAttendeeSQL)
	this.getEventAttendeesStmt, err = this.baseConn.Prepare(attendees.GetAttendeesByEventSQL)
	this.getUserEventsStmt, err = this.baseConn.Prepare(events.GetEventsForUserSQL)
	this.getEventStmt, err = this.baseConn.Prepare(events.GetEventSQL)
	this.getCoursesStmt, err = this.baseConn.Prepare(courses.GetCoursesSQL)
	this.getOptionsStmt, err = this.baseConn.Prepare(options.GetOptionsSQL)
	this.getSelectionStmt, err = this.baseConn.Prepare("SELECT `optionid` FROM `selections` WHERE `courseid` = ? AND `attendeeid` = ?")
	if err != nil {
		panic("Unable to prepare user statement")
	}
}

func (this *connection) GetUser(id int) (_ users.User, err error) {
	var rows *sql.Rows
	rows, err = this.getUserStmt.Query(id)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("User not found")
		return
	}
	return users.MakeUser(rows)
}

func (this *connection) GetUserByEmailPassword(email, password string) (user users.User, err error) {
	var rows *sql.Rows
	rows, err = this.getUserEmailStmt.Query(email)
	if err != nil {
		return
	}
	user, err = users.MakeUser(rows)
	if err != nil {
		return
	}
	err = user.VerifyPassword(password)
	return
}

func (this *connection) GetAttendee(id int) (_ attendees.Attendee, err error) {
	var rows *sql.Rows
	rows, err = this.getAttendeeStmt.Query(id)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("Attendee not found")
		return
	}
	return attendees.MakeAttendee(rows)
}

func (this *connection) GetAttendees(eventId int) (attendees []attendees.Attendee, err error) {
	var rows *sql.Rows
	attendees = make([]attendees.Attendee, 0, 5)
	rows, err = this.getEventAttendeesStmt.Query(eventId)
	if err != nil {
		return
	}
	var curId int
	var lastErr error
	var curAttendee attendees.Attendee
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&curId)
		curAttendee, lastErr = this.GetAttendee(curId)
		if lastErr == nil {
			attendees = append(attendees, curAttendee)
		}
	}
	if len(attendees) < 1 {
		err = lastErr
	}
	return
}

func (this *connection) GetAttendeeByKey(token string) (attendee attendees.Attendee, err error) {
	id, key := attendees.ParseToken(token)
	if (id == 0) || (key == "") {
		err = errors.New("Bad Request")
		return
	}
	attendee, err = this.GetAttendee(id)
	if err != nil {
		return
	}
	err = attendee.VerifyToken(key)
	return
}

func (this *connection) GetCourses(eventId int) (crses []courses.Course, err error) {
	var rows *sql.Rows
	crses = make([]courses.Course, 0, 5)
	rows, err = this.getCoursesStmt.Query(eventId)
	if err != nil {
		return
	}
	var curCourse courses.Course
	var lastErr error
	defer rows.Close()
	for rows.Next() {
		curCourse, lastErr = courses.MakeCourse(rows)
		if lastErr == nil {
			crses = append(crses, curCourse)
		}
	}
	if len(crses) < 1 {
		err = lastErr
	}
	return
}

func (this *connection) GetOptions(courseId int) (opts []options.Option, err error) {
	var rows *sql.Rows
	opts = make([]options.Option, 0, 5)
	rows, err = this.getOptionsStmt.Query(courseId)
	if err != nil {
		return
	}
	var curOpt options.Option
	var lastErr error
	defer rows.Close()
	for rows.Next() {
		curOpt, lastErr = options.MakeOption(rows)
		if lastErr == nil {
			opts = append(opts, curOpt)
		}
	}
	if len(opts) < 1 {
		err = lastErr
	}
	return
}

func (this *connection) GetEvent(id int) (_ events.Event, err error) {
	var rows *sql.Rows
	rows, err = this.getEventStmt.Query(id)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("Event not found")
		return
	}

	return events.MakeEvent(rows)
}

func (this *connection) GetEventsForUser(userId int) (events []events.Event, err error) {
	var rows *sql.Rows
	events = make([]events.Event, 0, 5)
	rows, err = this.getOptionsStmt.Query(userId)
	if err != nil {
		return
	}
	var curId int
	var lastErr error
	var curEvent events.Event
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&curId)
		curEvent, lastErr = this.GetEvent(curId)
		if lastErr == nil {
			events = append(events, curEvent)
		}
	}
	if len(events) < 1 {
		err = lastErr
	}
	return
}
