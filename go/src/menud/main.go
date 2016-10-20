package main

import (
	"goji.io"
	"goji.io/pat"
	"menud/routers/attendeeRouter"
	"menud/helpers/auth"
	"menud/config"
	"menud/database/connpool"
	"net/http"
	"menud/routers/adminRouter"
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

}
