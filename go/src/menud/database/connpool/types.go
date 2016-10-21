package connpool

import (
	"menud/components/attendees"
	"menud/components/courses"
	"menud/components/events"
	"menud/components/options"
	"menud/components/users"
)

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

type getCourseRequest struct {
	courseId int
	retChan  chan getCourseResponse
}
type getCourseResponse struct {
	crs courses.Course
	err error
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

type getOptionResponse struct {
	opt options.Option
	err error
}

type getOptionRequest struct {
	optionId int
	retChan  chan getOptionResponse
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
type setSelectionRequest struct {
	attendeeId  int
	courseId    int
	selectionId int
	retChan     chan getSelectionResponse
}
