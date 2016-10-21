package connpool

import "menud/config"

func init() {
	getUserChan = make(chan (getUserRequest))
	getUserByEmailPasswordChan = make(chan (getUserByEmailPasswordRequest))
	getAttendeeByKeyChan = make(chan (getAttendeeByKeyRequest))
	getAttendeeChan = make(chan (getAttendeeRequest))
	getAttendeesChan = make(chan (getAttendeesRequest))
	getCourseChan = make(chan (getCourseRequest))
	getCoursesChan = make(chan (getCoursesRequest))
	getOptionChan = make(chan (getOptionRequest))
	getOptionsChan = make(chan (getOptionsRequest))
	getSelectionChan = make(chan (getSelectionRequest))
	setSelectionChan = make(chan (setSelectionRequest))
	getEventChan = make(chan (getEventRequest))
	getEventsChan = make(chan (getEventsRequest))
	shutDownChan = make(chan (bool))

	for i := 0; i < config.ConnectionPoolSize(); i++ {
		go runListener()
	}
}

func runListener() {
	listener := &pooledConnection{}
	listener.setUp()
	listener.listen()
}

func Stop() {
	shutDownChan <- true
}

func Start() {
}
