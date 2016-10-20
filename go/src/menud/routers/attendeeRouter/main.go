package attendeeRouter

import (
	"encoding/json"
	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"
	"menud/components/attendees"
	"menud/components/courses"
	"menud/components/events"
	"menud/database/connpool"
	"menud/helpers/response"
	"net/http"
	"strconv"
)

func Get() *goji.Mux {
	rtr := goji.SubMux()
	rtr.HandleFuncC(pat.Get("/event"), getEvent)
	rtr.HandleFuncC(pat.Get("/courses"), getCourses)
	rtr.HandleFuncC(pat.Get("/options/:course"), getOptions)
	rtr.HandleFuncC(pat.Get("/selection/:course"), getSelection)
	rtr.HandleFuncC(pat.Post("/selection/:course"), setSelection)

	rtr.UseC(checkSession)
	return rtr
}

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

func getCourse(ctx context.Context, w http.ResponseWriter) (course courses.Course, bOK bool) {
	bOK = false
	event := ctx.Value(EventContextKey).(events.Event)
	courseIdStr := pat.Param(ctx, "course")
	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		response.ParseIdFailed(w)
		return
	}
	course, err = connpool.GetCourse(courseId)
	if (err != nil) || (course.EventId() != event.ID()) {
		response.CourseNotFound(w)
		return
	}
	bOK = true
	return
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
	selectionId, err = connpool.SetSelection(user.ID(), course.ID(), selectionId)
	if err != nil {
		response.InternalError(w, err)
		return
	}
	response.Send(w, selectionId)
}
