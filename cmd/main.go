package main

import (
	"flag"
	"net/http"

	"channeld.clewcat.com/channeld/pkg/channeld"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	/*
		getopt.Aliases(
			"sn", "serverNetwork",
			"sa", "serverAddress",
			"sfsm", "serverConnFSM",

			"cn", "clientNetwork",
			"ca", "clientAddress",
			"cfsm", "clientConnFSM",

			// "cs", "connSize",
		)
	*/

	sn := flag.String("sn", "tcp", "the network type for the server connections")
	sa := flag.String("sa", ":11288", "the network address for the server connections")
	sfsm := flag.String("sfsm", "../config/server_authoratative_fsm.json", "the path to the server FSM config")
	cn := flag.String("cn", "tcp", "the network type for the client connections")
	ca := flag.String("ca", ":12108", "the network address for the client connections")
	cfsm := flag.String("cfsm", "../config/client_non_authoratative_fsm.json", "the path to the client FSM config")
	// cs := flag.Int("cs", 1024, "the connection map buffer size")

	//getopt.Parse()
	flag.Parse()

	channeld.InitLogsAndMetrics()
	channeld.InitConnections(*sfsm, *cfsm)
	channeld.InitChannels()

	// Setup Prometheus
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8080", nil)

	go channeld.StartListening(channeld.SERVER, *sn, *sa)
	// FIXME: After all the server connections are established, the client connection should be listened.*/
	channeld.StartListening(channeld.CLIENT, *cn, *ca)

}
