# ReverseProxyDemo

# Sample 1: 

reverseproxy-transport 

Using the &httputil.ReverseProxy and define the transport which implements the RoundTripper interface

In this usecase end point validation proxy server (stop DELETE requests to hit main application)

# Sample 2:

reverseproxy-NewSingleHostReverseProxy

NewSingleHostReverseProxy will return the ReverseProxy and call the service using ServeHTTP method (ReverseProxy implements the ServeHTTP interface)

In this usecase rate limiter proxy server

In both the above cases, we have used CF Route service


https://docs.cloudfoundry.org/services/images/route-services-user-provided.png

Route service is what we have designed in demos.
