package auth

import (
"net/http"
"golang.org/x/net/context"
	"encoding/json"
	"menud/response"
	"menud/connpool"
	"menud/sessions"
	"menud/cookies"
)

type loginCreds struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func LoginUser(_ context.Context, w http.ResponseWriter, r *http.Request) {
	creds := loginCreds{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(creds)
	if err != nil {
		response.BadInput(w)
		return
	}
	user, err := connpool.GetUserByEmailPassword(creds.Email, creds.Password)
	if err != nil {
		response.BadLogin(w)
	}
	sess := sessions.New()
	sess.SetUserId(user.ID())

	http.SetCookie(w, cookies.Make(sess.GetSessionId()))
	response.Send(w, user)
}
