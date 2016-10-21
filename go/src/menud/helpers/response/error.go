package response

import (
	"fmt"
	"net/http"
)

const (
	errBadAttendeeToken  = iota
	errBadCredentials    = iota
	errBadInput          = iota
	errNeedAdminLogin    = iota
	errNeedAttendeeLogin = iota
	errInternal          = iota
	errParseIdFailed     = iota
	errCourseNotFound    = iota
	errSelectionNotFound = iota
	errOptionNotFound    = iota
	errEventNotFound     = iota
	errNotImplemented    = iota
)

var errorDict map[int]string = map[int]string{
	errBadAttendeeToken:  "Token not valid: please check you have copied the whole link from your email, or contact your event organiser",
	errBadCredentials:    "Email address or password incorrect.",
	errBadInput:          "Input JSON was not structured correctly",
	errNeedAdminLogin:    "You need to log in as an organiser to perform this action",
	errNeedAttendeeLogin: "You need to log in using your event invitation to perform this action",
	errInternal:          "An internal error occurred",
	errParseIdFailed:     "The requested ID could not be parsed",
	errCourseNotFound:    "The requested course could not be found",
	errSelectionNotFound: "You have not made a selection for the requested course",
	errOptionNotFound:    "The requested option is not available for the specified course",
	errNotImplemented:    "The requested function is not currently available",
	errEventNotFound:     "The requested event could not be found",
}

func BadToken(w http.ResponseWriter) {
	res := Response{}
	res.httpCode = 403
	res.ErrorCode = errBadAttendeeToken
	sendWithErrorMessage(res, w)
}
func BadInput(w http.ResponseWriter) {
	res := Response{}
	res.httpCode = 400
	res.ErrorCode = errBadInput
	sendWithErrorMessage(res, w)
}
func BadLogin(w http.ResponseWriter) {
	res := Response{}
	res.httpCode = 403
	res.ErrorCode = errBadCredentials
	sendWithErrorMessage(res, w)
}

func NeedAdminLogin(w http.ResponseWriter) {
	res := Response{}
	res.httpCode = 401
	res.ErrorCode = errNeedAdminLogin
	sendWithErrorMessage(res, w)
}
func NeedAttendeeLogin(w http.ResponseWriter) {
	res := Response{}
	res.httpCode = 401
	res.ErrorCode = errNeedAttendeeLogin
	sendWithErrorMessage(res, w)
}
func ParseIdFailed(w http.ResponseWriter) {
	res := Response{}
	res.httpCode = 400
	res.ErrorCode = errParseIdFailed
	sendWithErrorMessage(res, w)
}
func EventNotFound(w http.ResponseWriter) {
	res := Response{}
	res.httpCode = 404
	res.ErrorCode = errEventNotFound
	sendWithErrorMessage(res, w)
}
func CourseNotFound(w http.ResponseWriter) {
	res := Response{}
	res.httpCode = 404
	res.ErrorCode = errCourseNotFound
	sendWithErrorMessage(res, w)
}
func SelectionNotFound(w http.ResponseWriter) {
	res := Response{}
	res.httpCode = 404
	res.ErrorCode = errSelectionNotFound
	sendWithErrorMessage(res, w)
}
func OptionNotFound(w http.ResponseWriter) {
	res := Response{}
	res.httpCode = 404
	res.ErrorCode = errOptionNotFound
	sendWithErrorMessage(res, w)
}
func InternalError(w http.ResponseWriter, err error) {
	fmt.Println(err)
	res := Response{}
	res.httpCode = 500
	res.ErrorCode = errInternal
	sendWithErrorMessage(res, w)
}
func NotImplemented(w http.ResponseWriter) {
	res := Response{}
	res.httpCode = 500
	res.ErrorCode = errNotImplemented
	sendWithErrorMessage(res, w)
}

func sendWithErrorMessage(res Response, w http.ResponseWriter) {
	msg, ok := errorDict[res.ErrorCode]
	if !ok {
		msg = ""
	}
	res.ErrorMessage = msg
	writeJSON(res, w)
}
