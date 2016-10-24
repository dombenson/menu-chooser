package adminRouter

import (
	"golang.org/x/net/context"
	"menud/components/users"
	"menud/database/connpool"
	"menud/helpers/response"
	"net/http"
)

func getEvents(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	user := ctx.Value(UserContextKey).(users.User)
	ents, err := connpool.GetEventsForUser(user.ID())
	if err != nil {
		response.InternalError(w, err)
	}
	response.Send(w, ents)
}
