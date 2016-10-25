package db

import (
	_ "github.com/go-sql-driver/mysql"
	"menud/components/attendees"
	"menud/components/courses"
	"menud/components/events"
	"menud/components/options"
	"menud/components/users"
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

func GetConnection() Connection {
	obj := &connection{}
	obj.connect()
	return obj
}
