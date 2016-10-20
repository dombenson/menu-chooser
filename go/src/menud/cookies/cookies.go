package cookies

import (
	"net/http"
	"menud/config"
)

func Make(value string) *http.Cookie {
	return &http.Cookie{Name: config.CookieName(),
		Value: value, Secure: config.UseHTTPS(),
		HttpOnly: true, Path: config.PathPrefix()}
}