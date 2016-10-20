package attendeeRouter

import (
	"net/http"
	"goji.io"
	"golang.org/x/net/context"
	"menud/config"
	"menud/response"
	"menud/sessions"
	"menud/connpool"
)

type contextKey int

const (
	SessionContextKey contextKey = iota
	AttendeeContextKey contextKey = iota
	EventContextKey contextKey = iota
)

func checkSession (chain goji.Handler) goji.Handler {
	handler := &sessionChecker{chain}
	return handler
}

type sessionChecker struct {
	successHandler goji.Handler
}

func(this *sessionChecker ) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(config.CookieName())
	if err != nil {
		response.NeedAttendeeLogin(w)
		return
	}
	sess, err := sessions.Get(cookie.Value)
	if (err != nil) || (sess.GetAttendeeId() == 0) {
		response.NeedAttendeeLogin(w)
		return
	}
	user, err := connpool.GetAttendee(sess.GetAttendeeId())
	if err != nil {
		sessions.Destroy(sess.GetSessionId())
		response.NeedAttendeeLogin(w)
		return
	}
	if(user.EventId() != sess.GetEventId()) {
		sessions.Destroy(sess.GetSessionId())
		response.NeedAttendeeLogin(w)
		return
	}
	event, err := connpool.GetEvent(user.EventId())
	if err != nil {
		sessions.Destroy(sess.GetSessionId())
		response.NeedAttendeeLogin(w)
		return
	}
	childCtx := context.WithValue(context.WithValue(context.WithValue(ctx, AttendeeContextKey, user), SessionContextKey, sess), EventContextKey, event)
	this.successHandler.ServeHTTPC(childCtx, w, r)
}
