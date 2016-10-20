package auth

import (
	"net/http"
	"golang.org/x/net/context"
	"goji.io/pat"
	"menud/connpool"
	"menud/response"
	"menud/sessions"
	"menud/cookies"
)

func LoginAttendee(c context.Context, w http.ResponseWriter, r *http.Request) {
	token := pat.Param(c, "token")
	attendee, err := connpool.GetAttendeeByKey(token)
	if err != nil {
		response.BadToken(w)
		return
	}
	sess := sessions.New()
	sess.SetAttendeeId(attendee.ID())
	sess.SetEventId(attendee.EventId())

	http.SetCookie(w, cookies.Make(sess.GetSessionId()))
	response.Send(w, attendee)
}

