package main

import (
	"flag"
	"io"
	"log"
	"net/http"
)

type proxy struct {
}

func (p *proxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	log.Printf("remote_addr=%#v method=%#v url=%#v\n", req.RemoteAddr, req.Method, req.URL)

	resp, err := http.Get("https://pnpm.io" + req.URL.Path)
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

	handler := &proxy{}

	log.Println("Starting proxy server on", *addr)
	if err := http.ListenAndServe(*addr, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
