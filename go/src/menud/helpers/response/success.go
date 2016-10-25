package response

import "net/http"
import "menud/config"
import "fmt"

func Send(w http.ResponseWriter, data interface{}) {
	res := Response{}
	res.httpCode = 200
	res.ErrorCode = 0
	res.Data = data
	writeJSON(res, w)
}

func SendRedir(w http.ResponseWriter) {
        SendCorsHeaders(w)
        w.WriteHeader(200)

        strHtml := fmt.Sprintf(`<html><head><meta http-equiv="refresh" content="0; url=%s%s/"></head><body>Logging in</body></html>`, config.ExternalHost(), config.PathPrefix())
	w.Write([]byte(strHtml))
        return	
}
