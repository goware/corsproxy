package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var (
	flags = flag.NewFlagSet("corsproxy", flag.ExitOnError)

	fSource = flags.String("source", "", "Remote source to proxy to")
	fListen = flags.String("listen", "9090", "Local port to listen for this proxy service")
)

func main() {
	flags.Parse(os.Args[1:])

	if fSource == nil || *fSource == "" {
		fmt.Println("-source cannot be empty")
		os.Exit(1)
	}

	listen := *fListen
	if strings.Index(listen, ":") < 0 {
		listen = fmt.Sprintf("localhost:%s", listen)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/*", proxy())

	fmt.Printf("Proxying API requests from http://%s to %s ...\n", listen, *fSource)

	http.ListenAndServe(listen, r)
}

func proxy() http.Handler {
	source, err := url.Parse(*fSource)
	if err != nil {
		fmt.Printf("error parsing -source host %s because %v", *fSource, err)
		os.Exit(1)
	}

	director := func(req *http.Request) {
		req.URL.Scheme = source.Scheme
		req.URL.Host = source.Host
		req.Host = source.Host

		if req.Header.Get("Origin") != "" {
			req.Header.Set("Origin", source.String())
		}
	}

	modifyResponse := func(resp *http.Response) error {
		resp.Header.Set("Access-Control-Allow-Origin", "*")
		return nil
	}

	proxy := &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifyResponse,
	}

	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

	return proxy
}
