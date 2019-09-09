package main

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) < 1 {
		port = "8081"
	}
	mux := http.NewServeMux()
	mux.Handle("/", newProxy())
	http.ListenAndServe(":"+port, mux)

}
// return type is handler
func newProxy() http.Handler {
	proxy := &httputil.ReverseProxy{
		// here is the place to maintain the URL to call the main service
		Director: func(req *http.Request) {
			//cloud foundry set the header about the main app URL
			forwardedURL := req.Header.Get("X-Cf-Forwarded-Url")
			//you will get the request body , but you have set it again if you read otherwise it will be empty
			//i.e once read from req body , it will  get clear
			url, err := url.Parse(forwardedURL)
			if err != nil {
				log.Fatalln(err.Error())
			}

			req.URL = url
			req.Host = url.Host
		},
		Transport:      newSampleRoundTripper(), // assignment should be of http.RoundTripper type
	}
	return proxy
}

type SampleRoundTripper struct {
	transport http.RoundTripper //round tripper - if not maintained default transport will be used
}

func newSampleRoundTripper() *SampleRoundTripper {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &SampleRoundTripper{
		transport:   tr, // it can null and we dont need to maintain the transport
	}
}
//the transport is used to perform the proxy requests
// roundtrip interface should be called
func (r *SampleRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	println("Before calling the service")
	if req.Method == "DELETE"{
		resp := &http.Response{
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Not Accepted")),
		}
		return resp, nil
	}
	res, err := r.transport.RoundTrip(req)
	println("after calling the service , the response is " , res)
	return res, err
}
