package adminRouter

import (
	"encoding/csv"
	"goji.io/pat"
	"golang.org/x/net/context"
	"menud/components/events"
	"menud/components/options"
	"menud/database/connpool"
	"menud/helpers/email"
	"menud/helpers/response"
	"net/http"
	"strconv"
)

func getEventAttendees(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	event := ctx.Value(EventContextKey).(events.Event)
	attendees, err := connpool.GetAttendees(event.ID())
	if err != nil {
		response.InternalError(w, err)
		return
	}
	response.Send(w, attendees)
}
func getEventCourses(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	courses, err := connpool.GetCourses(ctx.Value(EventContextKey).(events.Event).ID())
	if err != nil {
		response.InternalError(w, err)
		return
	}
	response.Send(w, courses)
}
func sendInvites(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	event := ctx.Value(EventContextKey).(events.Event)
	attendees, err := connpool.GetAttendees(event.ID())
	if err != nil {
		response.InternalError(w, err)
		return
	}
	for _, att := range attendees {
		err = email.Send(att)
		if err != nil {
			response.InternalError(w, err)
			return
		}
	}
	response.Send(w, true)
}
func sendOneInvite(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	event := ctx.Value(EventContextKey).(events.Event)
	requestedAttId, err := strconv.Atoi(pat.Param(ctx, "attendeeId"))
	if err != nil {
		response.BadInput(w)
		return
	}
	requestedAtt, err := connpool.GetAttendee(requestedAttId)
	if err != nil {
		response.BadInput(w)
		return
	}
	if requestedAtt.EventId() != event.ID() {
		response.BadInput(w)
		return
	}
	err = email.Send(requestedAtt)
	if err != nil {
		response.InternalError(w, err)
		return
	}
	response.Send(w, true)
}
func newEvent(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func delEvent(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func setEventDetails(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func newAttendee(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func delAttendee(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func setAttendeeDetails(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func newCourse(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	response.NotImplemented(w)
}
func getEventCSV(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	event := ctx.Value(EventContextKey).(events.Event)
	courses, err := connpool.GetCourses(event.ID())
	if err != nil {
		response.InternalError(w, err)
		return
	}
	attendees, err := connpool.GetAttendees(event.ID())
	if err != nil {
		response.InternalError(w, err)
		return
	}
	colCount := 1 + len(courses)
	headerRow := make([]string, colCount)
	headerRow[0] = "Name"
	optInfo := make(map[int]options.Option)
	for crsIndex, course := range courses {
		headerRow[crsIndex+1] = course.Name()
		thisOpts, err := connpool.GetOptions(course.ID())
		if err != nil {
			response.InternalError(w, err)
			return
		}
		for _, opt := range thisOpts {
			optInfo[opt.ID()] = opt
		}
	}

	w.Header().Set("Content-Type", "text/csv")
	w.WriteHeader(200)

	csvWriter := csv.NewWriter(w)
	csvWriter.Write(headerRow)

	for _, attendee := range attendees {
		thisRow := make([]string, colCount)
		thisRow[0] = attendee.Name()
		for crsIndex, course := range courses {
			selName := ""
			curSel, err := connpool.GetSelection(attendee.ID(), course.ID())
			if err == nil {
				selOpt, ok := optInfo[curSel]
				if ok {
					selName = selOpt.Name()
				}
			}
			thisRow[crsIndex+1] = selName
		}
		csvWriter.Write(thisRow)
	}
	csvWriter.Flush()
}
