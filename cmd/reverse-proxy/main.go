package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var serverCount = 0

const (
	SERVER1 = "https://www.apple.com/"
	SERVER2 = "https://www.google.com/"
	SERVER3 = "https://www.facebook.com/"
	PORT    = "1338"
)

func main() {

	// start server
	http.HandleFunc("/", loadBalancer)
	log.Println("Listening for requests at http://localhost:8000/hello")
	log.Fatal(http.ListenAndServe(":"+PORT, nil)) //nil ... automatically creates http.NewServeMux() object
}

func loadBalancer(res http.ResponseWriter, req *http.Request) {
	// get address of one server
	url := getProxyURL()

	//log requets
	logRequestPayload(url)

	// forward request
	serveReverseProxy(url, res, req)

}

// get server using RR
func getProxyURL() string {
	serverPointer := 0
	var servers = []string{SERVER1, SERVER2, SERVER3}
	currentServer := servers[serverPointer]
	serverPointer++

	if serverPointer >= len(servers) {
		serverPointer = 0
	}

	return currentServer
}

// logger
func logRequestPayload(url string) {
	log.Println(url)
	file, err := os.OpenFile("logs/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	log.Println("Sending request to Proxy: " + url)
}

// create reverse proxy & serve http
func serveReverseProxy(targetUrl string, res http.ResponseWriter, req *http.Request) {
	// parse URL
	url, _ := url.Parse(targetUrl)

	// create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Serve (non-blocking; uses go-routine)
	proxy.ServeHTTP(res, req)
}
