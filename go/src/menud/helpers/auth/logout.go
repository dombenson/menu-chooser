package auth

import (
	"net/http"
	"menud/helpers/response"
	"menud/helpers/sessions"
	"menud/helpers/cookies"
	"menud/config"
	"golang.org/x/net/context"
	"time"
)

func Logout(_ context.Context, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(config.CookieName())
	if err != nil {
		response.Send(w, "Logged out")
		return
	}
	clearCookie := cookies.Make("")
	clearCookie.Expires = time.Now().Add(-1*time.Hour);
	http.SetCookie(w, clearCookie)

	sessions.Destroy(cookie.Value)

	response.Send(w, "Logged out")
}
