package attendeeRouter

import (
	"goji.io"
	"goji.io/pat"
)

func Get() *goji.Mux {
	rtr := goji.SubMux()
	rtr.HandleFuncC(pat.Get("/event"), getEvent)
	rtr.HandleFuncC(pat.Get("/courses"), getCourses)
	rtr.HandleFuncC(pat.Get("/options/:course"), getOptions)
	rtr.HandleFuncC(pat.Get("/selection/:course"), getSelection)
	rtr.HandleFuncC(pat.Post("/selection/:course"), setSelection)

	rtr.UseC(checkSession)
	return rtr
}
