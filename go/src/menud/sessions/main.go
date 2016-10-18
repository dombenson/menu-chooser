package sessions

type Session interface {
	GetUserId()(int)
	GetAttendeeId()(int)
	GetEventId()(int)
	SetUserId(int)()
	SetAttendeeId(int)()
	SetEventId(int)()
	GetSessionId()(string)
}

var sessionStoreLock sync.RWMutex
var sessionStore map[string]*session

func Init() {
	sessionStore = make(map[string]*session)
}

type session struct {
	userId int
	attendeeId int
	eventId int
	sessionId string
	expires time.Time
}

func(this *session) GetUserId() int {
	return this.userId
}
func(this *session) GetSessionId() string {
	return this.eventId
}
func(this *session) GetAttendeeId() int {
	return this.attendeeId
}
func(this *session) GetEventId() int {
	return this.sessionId
}

func(this *session) SetUserId(id int) {
	this.userId = id
}
func(this *session) SetAttendeeId(id int) {
	this.attendeeId = id
}
func(this *session) GetEventId(id int) {
	this.eventId = id
}

func Get(sessionId string)(sess Session, err error) {
	var bFound bool
	sessionStoreLock.RLock()
	sess, bFound = sessionStore[sessionId]
	sessionStoreLock.RUnLock()
	if !bFound {
		err = errors.New("Not found")
		return
	}
	if(sess.expires.After(time.Now())) {
		sess = nil
		err = errors.New("Expired")
		sessionStoreLock.Lock()
		delete(sessionStore[sessionId])
		sessionStoreLock.Unlock()
	} else {
		sess.expires = time.Now().Add(4*time.Hour)
	}
	return
}

func New()(sess Session) {
	bytes := make([]byte, 21)
	_, err := rand.Read(bytes)
	if err != nil {
		panic("Random source unavailable")
	}
	sessId := base64.StdEncoding.EncodeToString(bytes)
	newSess := &session{}
	newSess.sessionId = sessId
	newSess.expires = time.Now().Add(4*time.Hour)
	sessionStoreLock.Lock()
	sessionStore[sessId] = newSess
	sessionStoreLock.Unlock()
	return newSess
}

func Destroy(sessionId string) {
		sessionStoreLock.Lock()
		delete(sessionStore[sessionId])
		sessionStoreLock.Unlock()
}