package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

type Session interface {
	GetUserId() int
	GetAttendeeId() int
	GetEventId() int
	SetUserId(int)
	SetAttendeeId(int)
	SetEventId(int)
	GetSessionId() string
}

var sessionStoreLock sync.RWMutex
var sessionStore map[string]*session

func init() {
	sessionStore = make(map[string]*session)
}

type session struct {
	userId     int
	attendeeId int
	eventId    int
	sessionId  string
	expires    time.Time
}

func (this *session) GetUserId() int {
	return this.userId
}
func (this *session) GetSessionId() string {
	return this.sessionId
}
func (this *session) GetAttendeeId() int {
	return this.attendeeId
}
func (this *session) GetEventId() int {
	return this.eventId
}

func (this *session) SetUserId(id int) {
	this.userId = id
}
func (this *session) SetAttendeeId(id int) {
	this.attendeeId = id
}
func (this *session) SetEventId(id int) {
	this.eventId = id
}

func Get(sessionId string) (sess Session, err error) {
	var bFound bool
	sessionStoreLock.RLock()
	var locSess *session
	locSess, bFound = sessionStore[sessionId]
	sessionStoreLock.RUnlock()
	if !bFound {
		err = errors.New("Not found")
		return
	}
	if locSess.expires.After(time.Now()) {
		sess = nil
		err = errors.New("Expired")
		sessionStoreLock.Lock()
		delete(sessionStore, sessionId)
		sessionStoreLock.Unlock()
	} else {
		locSess.expires = time.Now().Add(4 * time.Hour)
		sess = locSess
	}
	return
}

func New() (sess Session) {
	bytes := make([]byte, 21)
	_, err := rand.Read(bytes)
	if err != nil {
		panic("Random source unavailable")
	}
	sessId := base64.StdEncoding.EncodeToString(bytes)
	newSess := &session{}
	newSess.sessionId = sessId
	newSess.expires = time.Now().Add(4 * time.Hour)
	sessionStoreLock.Lock()
	sessionStore[sessId] = newSess
	sessionStoreLock.Unlock()
	return newSess
}

func Destroy(sessionId string) {
	sessionStoreLock.Lock()
	delete(sessionStore, sessionId)
	sessionStoreLock.Unlock()
}
