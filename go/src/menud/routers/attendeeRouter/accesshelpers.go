package attendeeRouter

import (
	"goji.io/pat"
	"golang.org/x/net/context"
	"menud/components/courses"
	"menud/components/events"
	"menud/database/connpool"
	"menud/helpers/response"
	"net/http"
	"strconv"
)

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
