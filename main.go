package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
)

type proxy struct {
	remoteAddr string
}

func (p *proxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	log.Printf("remote_addr=%#v method=%#v url=%#v\n", req.RemoteAddr, req.Method, req.URL)

	resp, err := http.Get(p.remoteAddr + req.URL.Path)
	if err != nil {
		http.Error(wr, "Server Error", http.StatusInternalServerError)
		log.Println("ServeHTTP:", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("remote_addr=%#v status=%#v status_code=%#v\n", req.RemoteAddr, resp.Status, resp.StatusCode)

	for key, value := range resp.Header {
		wr.Header().Set(key, value[0])
	}
	wr.WriteHeader(resp.StatusCode)
	io.Copy(wr, resp.Body)
}

func main() {
	var addr = flag.String("addr", "0.0.0.0:8080", "The addr of the application.")
	flag.Parse()

	remoteAddr, ok := os.LookupEnv("REMOTE_ADDR")
	if !ok {
		log.Fatal("REMOTE_ADDR is not set")
	}

	handler := &proxy{
		remoteAddr: remoteAddr,
	}

	log.Printf("Starting proxy server on %s\n", *addr)
	if err := http.ListenAndServe(*addr, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
