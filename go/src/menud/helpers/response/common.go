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
	w.Header().Set("Content-Type", "application/json")
	if output.httpCode > 0 {
		w.WriteHeader(output.httpCode)
	}

	w.Write(marshalledJson)
	return
}
