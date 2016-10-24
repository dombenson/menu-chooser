package adminRouter

import (
	"goji.io"
	"goji.io/pat"
)

func Get() *goji.Mux {
	rtr := goji.SubMux()

	eventRtr := getEventRtr()
	courseRtr := getCourseRtr()

	rtr.HandleFuncC(pat.Get("/events"), getEvents)
	rtr.HandleFuncC(pat.Post("/events/new"), newEvent)

	rtr.HandleC(pat.New("/event/:eventId/*"), eventRtr)
	rtr.HandleC(pat.New("/course/:courseId/*"), courseRtr)

	rtr.UseC(checkSession)
	return rtr
}

func getEventRtr() *goji.Mux {

	eventRtr := goji.SubMux()
	eventRtr.HandleFuncC(pat.Get("/attendees"), getEventAttendees)
	eventRtr.HandleFuncC(pat.Get("/courses"), getEventCourses)

	eventRtr.HandleFuncC(pat.Post("/delete"), delEvent)
	eventRtr.HandleFuncC(pat.Post(""), setEventDetails)
	eventRtr.HandleFuncC(pat.Post("/attendee/new"), newAttendee)
	eventRtr.HandleFuncC(pat.Post("/attendee/:attendeeId/delete"), delAttendee)
	eventRtr.HandleFuncC(pat.Post("/attendee/:attendeeId"), setAttendeeDetails)
	eventRtr.HandleFuncC(pat.Post("/courses/new"), newCourse)

	eventRtr.UseC(checkEvent)

	return eventRtr
}

func getCourseRtr() *goji.Mux {

	courseRtr := goji.SubMux()

	courseRtr.HandleFuncC(pat.Get("/options"), getCourseOptions)
	courseRtr.HandleFuncC(pat.Get("/selections"), getCourseSelections)
	courseRtr.HandleFuncC(pat.Get("/"), getCourseDetails)

	courseRtr.HandleFuncC(pat.Post("/delete"), delCourse)
	courseRtr.HandleFuncC(pat.Post("/"), setCourseDetails)
	courseRtr.HandleFuncC(pat.Post("/option/new"), newOption)
	courseRtr.HandleFuncC(pat.Post("/option/:optionId/delete"), delOption)
	courseRtr.HandleFuncC(pat.Post("/option/:optionId"), setOptionDetails)

	courseRtr.UseC(checkCourse)

	return courseRtr
}
