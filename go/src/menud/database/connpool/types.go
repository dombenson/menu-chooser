package connpool

import (
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
	GetOptions(courseId int) (opts []options.Option, err error)
	GetEvent(id int) (_ events.Event, err error)
	GetEventsForUser(userId int) (events []events.Event, err error)
	GetSelection(attendeeId, courseId int) (optionId int, err error)
}

type getUserResponse struct {
	user users.User
	err  error
}

type getUserRequest struct {
	userId  int
	retChan chan getUserResponse
}

type getUserByEmailPasswordRequest struct {
	email    string
	password string
	retChan  chan getUserResponse
}

type getAttendeeResponse struct {
	attendee attendees.Attendee
	err      error
}

type getAttendeeRequest struct {
	attendeeId int
	retChan    chan getAttendeeResponse
}

type getAttendeeByKeyRequest struct {
	token   string
	retChan chan getAttendeeResponse
}

type getAttendeesResponse struct {
	attendees []attendees.Attendee
	err       error
}

type getAttendeesRequest struct {
	eventId int
	retChan chan getAttendeesResponse
}

type getCoursesResponse struct {
	crses []courses.Course
	err   error
}

type getCoursesRequest struct {
	eventId int
	retChan chan getCoursesResponse
}

type getOptionsResponse struct {
	opts []options.Option
	err  error
}

type getOptionsRequest struct {
	courseId int
	retChan  chan getOptionsResponse
}

type getEventResponse struct {
	event events.Event
	err   error
}

type getEventRequest struct {
	eventId int
	retChan chan getEventResponse
}

type getEventsResponse struct {
	event []events.Event
	err   error
}

type getEventsRequest struct {
	userId  int
	retChan chan getEventsResponse
}

type getSelectionResponse struct {
	optionId int
	err      error
}

type getSelectionRequest struct {
	attendeeId int
	courseId   int
	retChan    chan getSelectionResponse
}
