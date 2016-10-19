package connpool

import "menud/config"

func init() {
	getUserChan = make(chan (getUserRequest))
	getUserByEmailPasswordChan = make(chan (getUserByEmailPasswordRequest))
	getAttendeeByKeyChan = make(chan (getAttendeeByKeyRequest))
	getAttendeeChan = make(chan (getAttendeeRequest))
	getAttendeesChan = make(chan (getAttendeesRequest))
	getCoursesChan = make(chan (getCoursesRequest))
	getOptionsChan = make(chan (getOptionsRequest))
	getSelectionChan = make(chan (getSelectionRequest))
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
