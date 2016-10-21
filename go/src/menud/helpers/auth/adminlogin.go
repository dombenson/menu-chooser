package auth

import (
	"encoding/json"
	"golang.org/x/net/context"
	"menud/database/connpool"
	"menud/helpers/cookies"
	"menud/helpers/response"
	"menud/helpers/sessions"
	"net/http"
)

type loginCreds struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginUser(_ context.Context, w http.ResponseWriter, r *http.Request) {
	creds := loginCreds{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&creds)
	if err != nil {
		response.BadInput(w)
		return
	}
	user, err := connpool.GetUserByEmailPassword(creds.Email, creds.Password)
	if err != nil {
		response.BadLogin(w)
		return
	}
	sess := sessions.New()
	sess.SetUserId(user.ID())

	http.SetCookie(w, cookies.Make(sess.GetSessionId()))
	response.Send(w, user)
}
