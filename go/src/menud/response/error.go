package response

import "net/http"

const errBadAttendeeToken = 5;

var errorDict map[int]string = map[int]string{
	errBadAttendeeToken: "Token not valid: please check you have copied the whole link from your email, or contact your event organiser",
}


func BadToken(w http.ResponseWriter) {
	res := Response{}
	res.httpCode = 403
	res.ErrorCode = errBadAttendeeToken
	sendWithErrorMessage(res, w)
}

func sendWithErrorMessage(res Response, w http.ResponseWriter) {
	msg, ok := errorDict[res.ErrorCode]
	if(!ok) {
		msg = ""
	}
	res.ErrorMessage = msg
	writeJSON(res, w)
}