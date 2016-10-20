package main

import (
	"goji.io"
	"goji.io/pat"
	"menud/config"
	"menud/database/connpool"
	"menud/helpers/auth"
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
	topRouter.HandleFuncC(pat.Post("/adminlogin/"), auth.LoginUser)
	topRouter.HandleFuncC(pat.Get("/logout"), auth.Logout)

	topRouter.HandleC(pat.New("/admin/*"), admRouter)
	topRouter.HandleC(pat.New("/user/*"), attRouter)

	http.ListenAndServe(config.BindString(), topRouter)

	connpool.Stop()
}
