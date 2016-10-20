package response

import (
	"encoding/json"
	"net/http"
)

func writeJSON(output Response, w http.ResponseWriter) {
	marshalledJson, err := json.Marshal(output)
	if err != nil {
		http.Error(w, "Internal Error", 500)
		return
	}
	if output.httpCode > 0 {
		http.Error(w, "", output.httpCode)
	}
	w.Header().Add("Content-Type", "application/json")

	w.Write(marshalledJson)
	return
}
