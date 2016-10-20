package auth

import (
	"goji.io/pat"
	"golang.org/x/net/context"
	"menud/database/connpool"
	"menud/helpers/cookies"
	"menud/helpers/response"
	"menud/helpers/sessions"
	"net/http"
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
