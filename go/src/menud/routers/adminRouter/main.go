package adminRouter

import (
	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"
	"menud/components/users"
	"menud/database/connpool"
	"menud/helpers/response"
	"net/http"
)

func Get() *goji.Mux {
	rtr := goji.SubMux()
	rtr.HandleFuncC(pat.Get("/events"), getEvents)

	rtr.UseC(checkSession)
	return rtr
}

func getEvents(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	user := ctx.Value(UserContextKey).(users.User)
	ents, err := connpool.GetEventsForUser(user.ID())
	if err != nil {
		response.InternalError(w, err)
	}
	response.Send(w, ents)
}
