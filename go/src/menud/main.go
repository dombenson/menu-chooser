package main

import (
	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"
	"menud/config"
	"menud/database/connpool"
	"menud/helpers/auth"
	"menud/helpers/response"
	"menud/routers/adminRouter"
	"menud/routers/attendeeRouter"
	"net/http"
)

func main() {
	connpool.Start()

	topRouter := goji.NewMux()

	admRouter := adminRouter.Get()

	attRouter := attendeeRouter.Get()

	topRouter.HandleFuncC(pat.Get("/login/:token"), auth.LoginAttendee)
	topRouter.HandleFuncC(pat.Post("/login"), auth.LoginAttendeePost)
	topRouter.HandleFuncC(pat.Options("/*"), sendCors)
	topRouter.HandleFuncC(pat.Post("/adminlogin"), auth.LoginUser)
	topRouter.HandleFuncC(pat.Get("/logout"), auth.Logout)

	topRouter.HandleC(pat.New("/admin/*"), admRouter)
	topRouter.HandleC(pat.New("/user/*"), attRouter)

	http.ListenAndServe(config.BindString(), topRouter)

	connpool.Stop()
}

func sendCors(_ context.Context, w http.ResponseWriter, _ *http.Request) {
	response.SendCorsHeaders(w)
}
