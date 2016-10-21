package db

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"menud/components/attendees"
	"menud/components/courses"
	"menud/components/events"
	"menud/components/options"
	"menud/components/users"
	"menud/config"
)

type Connection interface {
	GetUser(int) (users.User, error)
	GetUserByEmailPassword(email, password string) (users.User, error)
	GetAttendee(id int) (_ attendees.Attendee, err error)
	GetAttendees(eventId int) (attendees []attendees.Attendee, err error)
	GetAttendeeByKey(token string) (attendee attendees.Attendee, err error)
	GetCourses(eventId int) (crses []courses.Course, err error)
	GetCourse(courseId int) (crs courses.Course, err error)
	GetOption(optionID int) (opt options.Option, err error)
	GetOptions(courseId int) (opts []options.Option, err error)
	GetEvent(id int) (_ events.Event, err error)
	GetEventsForUser(userId int) (events []events.Event, err error)
	GetSelection(attendeeId, courseId int) (optionId int, err error)
	SetSelection(attendeeId, courseId, selectionId int) (optionId int, err error)
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
	getCourseStmt         *sql.Stmt
	getOptionStmt         *sql.Stmt
	getOptionsStmt        *sql.Stmt
	getSelectionStmt      *sql.Stmt
	setSelectionStmt      *sql.Stmt
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
	if err != nil {
		panic("Unable to prepare user statement: " + err.Error())
	}
	this.getUserEmailStmt, err = this.baseConn.Prepare(users.GetUserByEmailSQL)
	if err != nil {
		panic("Unable to prepare get user by email statement: " + err.Error())
	}
	this.getAttendeeStmt, err = this.baseConn.Prepare(attendees.GetAttendeeSQL)
	if err != nil {
		panic("Unable to prepare get attendee statement: " + err.Error())
	}
	this.getEventAttendeesStmt, err = this.baseConn.Prepare(attendees.GetAttendeesByEventSQL)
	if err != nil {
		panic("Unable to prepare get event attendees statement: " + err.Error())
	}
	this.getUserEventsStmt, err = this.baseConn.Prepare(events.GetEventsForUserSQL)
	if err != nil {
		panic("Unable to prepare get events for user statement: " + err.Error())
	}
	this.getEventStmt, err = this.baseConn.Prepare(events.GetEventSQL)
	if err != nil {
		panic("Unable to prepare get event statement: " + err.Error())
	}
	this.getCoursesStmt, err = this.baseConn.Prepare(courses.GetCoursesSQL)
	if err != nil {
		panic("Unable to prepare get courses statement: " + err.Error())
	}
	this.getCourseStmt, err = this.baseConn.Prepare(courses.GetCourseSQL)
	if err != nil {
		panic("Unable to prepare get course statement: " + err.Error())
	}
	this.getOptionStmt, err = this.baseConn.Prepare(options.GetOptionSQL)
	if err != nil {
		panic("Unable to prepare get option statement: " + err.Error())
	}
	this.getOptionsStmt, err = this.baseConn.Prepare(options.GetOptionsSQL)
	if err != nil {
		panic("Unable to prepare get options for course statement: " + err.Error())
	}
	this.getSelectionStmt, err = this.baseConn.Prepare("SELECT `optionid` FROM `selections` WHERE `courseid` = ? AND `attendeeid` = ?")
	if err != nil {
		panic("Unable to prepare get selection statement: " + err.Error())
	}
	this.setSelectionStmt, err = this.baseConn.Prepare("INSERT INTO `selections` (`courseid`,`attendeeid`,`optionid`) VALUES (?,?,?) ON DUPLICATE KEY UPDATE `optionid` = VALUES(`optionid`)")
	if err != nil {
		panic("Unable to prepare get selection statement: " + err.Error())
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
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("User not found")
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

func (this *connection) GetAttendees(eventId int) (atts []attendees.Attendee, err error) {
	var rows *sql.Rows
	atts = make([]attendees.Attendee, 0, 5)
	rows, err = this.getEventAttendeesStmt.Query(eventId)
	if err != nil {
		return
	}
	var curId int
	var lastErr error
	var curAttendee attendees.Attendee
	defer rows.Close()
	for rows.Next() {
		lastErr = rows.Scan(&curId)
		if lastErr != nil {
			continue
		}
		curAttendee, lastErr = this.GetAttendee(curId)
		if lastErr == nil {
			atts = append(atts, curAttendee)
		}
	}
	if len(atts) < 1 {
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

func (this *connection) GetCourse(courseId int) (_ courses.Course, err error) {
	var rows *sql.Rows
	rows, err = this.getCourseStmt.Query(courseId)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("Attendee not found")
		return
	}
	return courses.MakeCourse(rows)
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

func (this *connection) GetOption(id int) (_ options.Option, err error) {
	var rows *sql.Rows
	rows, err = this.getOptionStmt.Query(id)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("Event not found")
		return
	}

	return options.MakeOption(rows)
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

func (this *connection) GetEventsForUser(userId int) (ents []events.Event, err error) {
	var rows *sql.Rows
	ents = make([]events.Event, 0, 5)
	rows, err = this.getUserEventsStmt.Query(userId)
	if err != nil {
		return
	}
	var curId int
	var lastErr error
	var curEvent events.Event
	defer rows.Close()
	for rows.Next() {
		lastErr = rows.Scan(&curId)
		if lastErr == nil {
			curEvent, lastErr = this.GetEvent(curId)
			if lastErr == nil {
				ents = append(ents, curEvent)
			}
		}
	}
	if len(ents) < 1 {
		err = lastErr
	}
	return
}

func (this *connection) GetSelection(attendeeId, courseId int) (optionId int, err error) {
	var rows *sql.Rows
	rows, err = this.getSelectionStmt.Query(courseId, attendeeId)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("Event not found")
		return
	}
	err = rows.Scan(&optionId)
	return
}

func (this *connection) SetSelection(attendeeId, courseId, optionId int) (outOptionId int, err error) {
	outOptionId, _ = this.GetSelection(attendeeId, courseId)
	_, err = this.setSelectionStmt.Exec(courseId, attendeeId, optionId)
	if err != nil {
		return
	}
	outOptionId = optionId
	return
}
