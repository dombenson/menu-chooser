package connpool

var shutDownChan chan (bool)

var getUserChan chan (getUserRequest)
var getUserByEmailPasswordChan chan (getUserByEmailPasswordRequest)
var getAttendeeByKeyChan chan (getAttendeeByKeyRequest)
var getAttendeeChan chan (getAttendeeRequest)
var getAttendeesChan chan (getAttendeesRequest)
var getCoursesChan chan (getCoursesRequest)
var getOptionsChan chan (getOptionsRequest)
var getSelectionChan chan (getSelectionRequest)
var getEventChan chan (getEventRequest)
var getEventsChan chan (getEventsRequest)
