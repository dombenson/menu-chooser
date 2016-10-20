package response

import "net/http"

func Send(w http.ResponseWriter, data interface{}) {
	res := Response{}
	res.httpCode = 200
	res.ErrorCode = 0
	res.Data = data
	writeJSON(res, w)
}

