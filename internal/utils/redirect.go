package utils

import (
	"net/http"
)

// Redirect function that redirects HTTP to HTTPS, taking into account the host and port
func Redirect(listenAddress string, httpsPort string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		target := "https://" + listenAddress
		if httpsPort != "443" {
			target += ":" + httpsPort
		}
		target += req.URL.Path
		http.Redirect(w, req, target, http.StatusMovedPermanently)
	}
}
