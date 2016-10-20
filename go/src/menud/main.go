package main

import (
	"menud/connpool"
	"goji.io"
	"goji.io/pat"
	"net/http"
	"menud/config"
	"menud/auth"
	"menud/attendeeRouter"
	"menud/adminRouter"
)

func main() {
	// Force the connection pool to start up so that statements are prepared and we can bail out immediately if
	// there is a problem with them (rather than panic-ing on a user call)
	connpool.GetUser(1)

	topRouter := goji.NewMux()

	admRouter := adminRouter.Get()

	attRouter := attendeeRouter.Get()

	topRouter.HandleFuncC(pat.Get("/login/:token"), auth.LoginAttendee)
	topRouter.HandleFuncC(pat.Post("/adminlogin/"), auth.LoginUser)

	topRouter.HandleC(pat.New("/admin/*"), admRouter)
	topRouter.HandleC(pat.New("/user/*"), attRouter)

	http.ListenAndServe(config.BindString(), topRouter)

}