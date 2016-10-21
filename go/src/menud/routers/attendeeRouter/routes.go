package attendeeRouter

import (
	"encoding/json"
	"golang.org/x/net/context"
	"menud/components/attendees"
	"menud/components/events"
	"menud/database/connpool"
	"menud/helpers/response"
	"net/http"
)

func getEvent(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	event := ctx.Value(EventContextKey).(events.Event)
	response.Send(w, event)
}

func getCourses(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	event := ctx.Value(EventContextKey).(events.Event)
	crses, err := connpool.GetCourses(event.ID())
	if err != nil {
		response.InternalError(w, err)
		return
	}
	response.Send(w, crses)
}

func getOptions(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	course, ok := getCourse(ctx, w)
	if !ok {
		return
	}
	opts, err := connpool.GetOptions(course.ID())
	if err != nil {
		response.CourseNotFound(w)
		return
	}
	response.Send(w, opts)
}

func getSelection(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := ctx.Value(AttendeeContextKey).(attendees.Attendee)
	course, ok := getCourse(ctx, w)
	if !ok {
		return
	}
	selectionId, err := connpool.GetSelection(user.ID(), course.ID())
	if err != nil {
		response.SelectionNotFound(w)
		return
	}
	response.Send(w, selectionId)
}

func setSelection(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	var selectionId int
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&selectionId)
	user := ctx.Value(AttendeeContextKey).(attendees.Attendee)
	course, ok := getCourse(ctx, w)
	if !ok {
		return
	}
	option, err := connpool.GetOption(selectionId)
	if err != nil {
		response.OptionNotFound(w)
		return
	}
	if option.CourseID() != course.ID() {
		response.OptionNotFound(w)
		return
	}
	selectionId, err = connpool.SetSelection(user.ID(), course.ID(), selectionId)
	if err != nil {
		response.InternalError(w, err)
		return
	}
	response.Send(w, selectionId)
}
