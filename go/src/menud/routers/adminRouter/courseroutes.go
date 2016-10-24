package adminRouter

import (
	"golang.org/x/net/context"
	"menud/components/courses"
	"menud/components/events"
	"menud/database/connpool"
	"menud/helpers/response"
	"net/http"
)

func getCourseDetails(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.Send(w, ctx.Value(CourseContextKey).(courses.Course))
}
func getCourseOptions(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	opts, err := connpool.GetOptions(ctx.Value(CourseContextKey).(courses.Course).ID())
	if err != nil {
		response.InternalError(w, err)
		return
	}
	response.Send(w, opts)
}
func getCourseSelections(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	ret := make(map[int]int)
	course := ctx.Value(CourseContextKey).(courses.Course)
	atts, err := connpool.GetAttendees(ctx.Value(EventContextKey).(events.Event).ID())
	if err != nil {
		response.InternalError(w, err)
		return
	}
	for _, att := range atts {
		thisSel, err := connpool.GetSelection(att.ID(), course.ID())
		if err != nil {
			response.InternalError(w, err)
			return
		}
		ret[att.ID()] = thisSel
	}
	response.Send(w, ret)
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
