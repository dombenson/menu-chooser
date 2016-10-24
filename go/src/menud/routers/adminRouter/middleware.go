package adminRouter

import (
	"goji.io"
)

type contextKey int

const (
	SessionContextKey contextKey = iota
	UserContextKey    contextKey = iota
	EventContextKey   contextKey = iota
	CourseContextKey  contextKey = iota
)

func checkSession(chain goji.Handler) goji.Handler {
	handler := &sessionChecker{chain}
	return handler
}

func checkEvent(chain goji.Handler) goji.Handler {
	handler := &eventChecker{chain}
	return handler
}

func checkCourse(chain goji.Handler) goji.Handler {
	handler := &courseChecker{chain}
	return handler
}
