package adminRouter

import (
	"golang.org/x/net/context"
	"menud/components/events"
	"menud/database/connpool"
	"menud/helpers/response"
	"net/http"
)

func getEventAttendees(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	event := ctx.Value(EventContextKey).(events.Event)
	attendees, err := connpool.GetAttendees(event.ID())
	if err != nil {
		response.InternalError(w, err)
	}
	response.Send(w, attendees)
}
func getEventCourses(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	courses, err := connpool.GetCourses(ctx.Value(EventContextKey).(events.Event).ID())
	if err != nil {
		response.InternalError(w, err)
	}
	response.Send(w, courses)
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
