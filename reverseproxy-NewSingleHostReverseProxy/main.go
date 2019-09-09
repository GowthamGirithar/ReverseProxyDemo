package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

import "golang.org/x/time/rate"

var limiter = rate.NewLimiter(2, 5) // r tokens will be put in bucket every second and 5 is max size of bucket
func main() {
	port := os.Getenv("PORT")
	if len(port) < 1 {
		port = "8081"
	}
	mux := http.NewServeMux()
	HandlerFunction(mux)
	http.ListenAndServe(":"+port, mux)

}
func HandlerFunction(mux *http.ServeMux) {
	mux.Handle("/", RateLimiter())
}

func RateLimiter() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if limiter.Allow() {
			//cloud foundry set the header about the main app URL
			forwardedURL := request.Header.Get("X-Cf-Forwarded-Url")
			// IT WILL CALL THE SERVICE AND THEN RETURN THE RESPONSE FROM THER
			URL, _ := url.Parse(forwardedURL)
			// RETURN VALUE OF NewSingleHostReverseProxy IS ReverseProxy and it also implement ServeHTTP
			httputil.NewSingleHostReverseProxy(URL).ServeHTTP(writer, request)
			//normally we use next.serverhttp if we have to call other middleware in same service
			//for calling other services , reverseproxy concept should be used
		} else {
			// send as too many requests
			http.Error(writer, http.StatusText(429), http.StatusTooManyRequests)
		}

	})
}
