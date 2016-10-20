package adminRouter

import (
	"net/http"
	"goji.io"
	"menud/config"
	"menud/response"
	"menud/sessions"
	"menud/connpool"
	"golang.org/x/net/context"
)

type contextKey int

const (
	SessionContextKey contextKey = iota
	UserContextKey contextKey = iota
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
