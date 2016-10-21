package adminRouter

import (
	"goji.io"
	"goji.io/pat"
)

func Get() *goji.Mux {
	rtr := goji.SubMux()

	eventRtr := goji.SubMux()
	courseRtr := goji.SubMux()

	rtr.HandleFuncC(pat.Get("/events"), getEvents)
	rtr.HandleFuncC(pat.Post("/events/new"), newEvent)

	rtr.HandleC(pat.New("/event/:eventId/*"), eventRtr)
	rtr.HandleC(pat.New("/course/:courseId/*"), courseRtr)

	eventRtr.HandleFuncC(pat.Get("/attendees"), getEventAttendees)
	eventRtr.HandleFuncC(pat.Get("/courses"), getEventCourses)

	courseRtr.HandleFuncC(pat.Get("/options"), getCourseOptions)
	courseRtr.HandleFuncC(pat.Get("/selections"), getCourseSelections)

	eventRtr.HandleFuncC(pat.Post("/delete"), delEvent)
	eventRtr.HandleFuncC(pat.Post(""), setEventDetails)
	eventRtr.HandleFuncC(pat.Post("/attendee/new"), newAttendee)
	eventRtr.HandleFuncC(pat.Post("/attendee/:attendeeId/delete"), delAttendee)
	eventRtr.HandleFuncC(pat.Post("/attendee/:attendeeId"), setAttendeeDetails)
	eventRtr.HandleFuncC(pat.Post("/courses/new"), newCourse)

	courseRtr.HandleFuncC(pat.Post("/delete"), delCourse)
	courseRtr.HandleFuncC(pat.Post("/"), setCourseDetails)
	courseRtr.HandleFuncC(pat.Post("/option/new"), newOption)
	courseRtr.HandleFuncC(pat.Post("/option/:optionId/delete"), delOption)
	courseRtr.HandleFuncC(pat.Post("/option/:optionId"), setOptionDetails)

	eventRtr.UseC(checkEvent)

	rtr.UseC(checkSession)
	return rtr
}
