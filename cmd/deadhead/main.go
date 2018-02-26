package main

import (
	"flag"
	"net"

	"github.com/netvm/funcs/helloworld"
	"github.com/sirupsen/logrus"

	"github.com/netvm/netvm"
)

var addr = flag.String("addr", ":9905", "Host address to listen at")

func main() {
	flag.Parse()

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		logrus.WithError(err).Fatal("Unable to start listener")
	}

	onError := func(err error) {
		logrus.WithError(err).Fatal("Error in deadhead")
	}

	err = netvm.ServeDeadhead(
		l,
		netvm.FuncHydrater{
			"90210": helloworld.ServeHTTP,
		},
		onError,
	)

	if err != nil {
		logrus.WithError(err).Fatal("Unable to start server")
	}
}
