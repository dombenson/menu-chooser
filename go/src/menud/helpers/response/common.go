package response

import (
	"encoding/json"
	"menud/config"
	"net/http"
)

func writeJSON(output Response, w http.ResponseWriter) {
	marshalledJson, err := json.Marshal(output)
	if err != nil {
		http.Error(w, "Internal Error", 500)
		return
	}
	SendCorsHeaders(w)
	if output.httpCode > 0 {
		w.WriteHeader(output.httpCode)
	}

	w.Write(marshalledJson)
	return
}

func SendCorsHeaders(w http.ResponseWriter) {
	if config.CorsEnabled() {
		w.Header().Set("Access-Control-Allow-Origin", config.CorsOrigin())
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
	}
}
