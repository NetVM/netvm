package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/netvm/netvm"
	"github.com/sirupsen/logrus"
)

var hostMap = map[string]string{
	"func1.netvmnetworknet.com": "90210",
}

var addr = flag.String("addr", ":9906", "Host address to listen at")
var deadheadSvc = flag.String("deadhead", "localhost:9905", "Host to proxy requests to")

func main() {
	flag.Parse()

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		logrus.WithError(err).Fatal("Unable to start listener")
	}

	target := &url.URL{
		Scheme: "http",
		Host:   *deadheadSvc,
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Director = func(r *http.Request) {
		r.URL.Host = target.Host
		r.URL.Scheme = target.Scheme
		r.URL.Path = target.Path

		sourceHost := strings.Split(r.Host, ":")[0]
		id := hostMap[sourceHost]

		r.Header.Set(netvm.HydrationIDHeader, id)
	}

	log.Fatal(http.Serve(l, proxy))
}
