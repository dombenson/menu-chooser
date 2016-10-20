package cookies

import (
	"menud/config"
	"net/http"
)

func Make(value string) *http.Cookie {
	return &http.Cookie{Name: config.CookieName(),
		Value: value, Secure: config.UseHTTPS(),
		HttpOnly: true, Path: config.PathPrefix()}
}
