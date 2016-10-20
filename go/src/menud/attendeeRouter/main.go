package attendeeRouter

import "goji.io"

func Get() (*goji.Mux) {
	rtr := goji.SubMux()
	rtr.UseC(checkSession)
	return rtr
}

