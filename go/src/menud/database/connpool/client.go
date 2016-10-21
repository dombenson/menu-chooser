package connpool

import (
	"menud/components/attendees"
	"menud/components/courses"
	"menud/components/events"
	"menud/components/options"
	"menud/components/users"
)

func GetUser(userId int) (users.User, error) {
	req := getUserRequest{}
	req.retChan = make(chan (getUserResponse))
	req.userId = userId
	getUserChan <- req
	res := <-req.retChan
	return res.user, res.err
}
func GetUserByEmailPassword(email, password string) (users.User, error) {
	req := getUserByEmailPasswordRequest{}
	req.retChan = make(chan (getUserResponse))
	req.email = email
	req.password = password
	getUserByEmailPasswordChan <- req
	res := <-req.retChan
	return res.user, res.err
}
func GetAttendee(id int) (attendees.Attendee, error) {
	req := getAttendeeRequest{}
	req.retChan = make(chan (getAttendeeResponse))
	req.attendeeId = id
	getAttendeeChan <- req
	res := <-req.retChan
	return res.attendee, res.err
}
func GetAttendees(eventId int) ([]attendees.Attendee, error) {
	req := getAttendeesRequest{}
	req.retChan = make(chan (getAttendeesResponse))
	req.eventId = eventId
	getAttendeesChan <- req
	res := <-req.retChan
	return res.attendees, res.err
}
func GetAttendeeByKey(token string) (attendees.Attendee, error) {
	req := getAttendeeByKeyRequest{}
	req.retChan = make(chan (getAttendeeResponse))
	req.token = token
	getAttendeeByKeyChan <- req
	res := <-req.retChan
	return res.attendee, res.err
}
func GetCourse(courseId int) (courses.Course, error) {
	req := getCourseRequest{}
	req.retChan = make(chan (getCourseResponse))
	req.courseId = courseId
	getCourseChan <- req
	res := <-req.retChan
	return res.crs, res.err
}
func GetCourses(eventId int) ([]courses.Course, error) {
	req := getCoursesRequest{}
	req.retChan = make(chan (getCoursesResponse))
	req.eventId = eventId
	getCoursesChan <- req
	res := <-req.retChan
	return res.crses, res.err
}
func GetOptions(courseId int) ([]options.Option, error) {
	req := getOptionsRequest{}
	req.retChan = make(chan (getOptionsResponse))
	req.courseId = courseId
	getOptionsChan <- req
	res := <-req.retChan
	return res.opts, res.err
}
func GetOption(optionId int) (options.Option, error) {
	req := getOptionRequest{}
	req.retChan = make(chan (getOptionResponse))
	req.optionId = optionId
	getOptionChan <- req
	res := <-req.retChan
	return res.opt, res.err
}
func GetEvent(id int) (events.Event, error) {
	req := getEventRequest{}
	req.retChan = make(chan (getEventResponse))
	req.eventId = id
	getEventChan <- req
	res := <-req.retChan
	return res.event, res.err
}
func GetEventsForUser(userId int) ([]events.Event, error) {
	req := getEventsRequest{}
	req.retChan = make(chan (getEventsResponse))
	req.userId = userId
	getEventsChan <- req
	res := <-req.retChan
	return res.event, res.err
}
func GetSelection(attendeeId, courseId int) (int, error) {
	req := getSelectionRequest{}
	req.retChan = make(chan (getSelectionResponse))
	req.attendeeId = attendeeId
	req.courseId = courseId
	getSelectionChan <- req
	res := <-req.retChan
	return res.optionId, res.err
}

func SetSelection(attendeeId, courseId, selectionId int) (int, error) {
	req := setSelectionRequest{}
	req.retChan = make(chan (getSelectionResponse))
	req.attendeeId = attendeeId
	req.courseId = courseId
	req.selectionId = selectionId
	setSelectionChan <- req
	res := <-req.retChan
	return res.optionId, res.err
}
