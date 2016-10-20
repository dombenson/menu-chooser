package auth

import (
	"golang.org/x/net/context"
	"menud/config"
	"menud/helpers/cookies"
	"menud/helpers/response"
	"menud/helpers/sessions"
	"net/http"
	"time"
)

func Logout(_ context.Context, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(config.CookieName())
	if err != nil {
		response.Send(w, "Logged out")
		return
	}
	clearCookie := cookies.Make("")
	clearCookie.Expires = time.Now().Add(-1 * time.Hour)
	http.SetCookie(w, clearCookie)

	sessions.Destroy(cookie.Value)

	response.Send(w, "Logged out")
}
