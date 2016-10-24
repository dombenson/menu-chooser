package adminRouter

import (
	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"
	"menud/components/users"
	"menud/config"
	"menud/database/connpool"
	"menud/helpers/response"
	"menud/helpers/sessions"
	"net/http"
	"strconv"
)

type sessionChecker struct {
	successHandler goji.Handler
}

func (this *sessionChecker) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(config.CookieName())
	if err != nil {
		response.NeedAdminLogin(w)
		return
	}
	sess, err := sessions.Get(cookie.Value)
	if (err != nil) || (sess.GetUserId() == 0) {
		response.NeedAdminLogin(w)
		return
	}
	user, err := connpool.GetUser(sess.GetUserId())
	if err != nil {
		response.NeedAdminLogin(w)
		return
	}
	childCtx := context.WithValue(context.WithValue(ctx, UserContextKey, user), SessionContextKey, sess)
	this.successHandler.ServeHTTPC(childCtx, w, r)
}

type eventChecker struct {
	successHandler goji.Handler
}

func (this *eventChecker) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	requestedEventId, err := strconv.Atoi(pat.Param(ctx, "eventId"))
	if err != nil || requestedEventId == 0 {
		response.BadInput(w)
		return
	}
	eventDetails, err := connpool.GetEvent(requestedEventId)
	if err != nil {
		response.EventNotFound(w)
		return
	}
	userDetails := ctx.Value(UserContextKey).(users.User)
	if eventDetails.UserID() != userDetails.ID() {

		response.EventNotFound(w)
		return
	}

	childCtx := context.WithValue(ctx, EventContextKey, eventDetails)
	this.successHandler.ServeHTTPC(childCtx, w, r)
}

type courseChecker struct {
	successHandler goji.Handler
}

func (this *courseChecker) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	requestedCourseId, err := strconv.Atoi(pat.Param(ctx, "courseId"))
	if err != nil || requestedCourseId == 0 {
		response.BadInput(w)
		return
	}
	courseDetails, err := connpool.GetCourse(requestedCourseId)
	if err != nil {
		response.CourseNotFound(w)
		return
	}
	eventDetails, err := connpool.GetEvent(courseDetails.EventId())
	if err != nil {
		response.CourseNotFound(w)
		return
	}
	userDetails := ctx.Value(UserContextKey).(users.User)
	if eventDetails.UserID() != userDetails.ID() {

		response.CourseNotFound(w)
		return
	}

	childCtx := context.WithValue(context.WithValue(ctx, EventContextKey, eventDetails), CourseContextKey, courseDetails)
	this.successHandler.ServeHTTPC(childCtx, w, r)
}
