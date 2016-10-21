package adminRouter

import (
	"golang.org/x/net/context"
	"menud/components/events"
	"menud/components/users"
	"menud/database/connpool"
	"menud/helpers/response"
	"net/http"
)

func getEvents(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	user := ctx.Value(UserContextKey).(users.User)
	ents, err := connpool.GetEventsForUser(user.ID())
	if err != nil {
		response.InternalError(w, err)
	}
	response.Send(w, ents)
}

func getEventAttendees(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	event := ctx.Value(EventContextKey).(events.Event)
	attendees, err := connpool.GetAttendees(event.ID())
	if err != nil {
		response.InternalError(w, err)
	}
	response.Send(w, attendees)
}
func getEventCourses(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func getCourseOptions(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func getCourseSelections(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func newEvent(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func delEvent(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func setEventDetails(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func newAttendee(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func delAttendee(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func setAttendeeDetails(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func newCourse(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func delCourse(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func setCourseDetails(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func newOption(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func delOption(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func setOptionDetails(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
