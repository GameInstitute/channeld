package main

import (
	"fmt"
	"net/http"

	"channeld.clewcat.com/channeld/pkg/channeld"
	"channeld.clewcat.com/channeld/pkg/channeldpb"
	"channeld.clewcat.com/channeld/pkg/unreal"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	if err := channeld.GlobalSettings.ParseFlag(); err != nil {
		fmt.Printf("error parsing CLI flag: %v\n", err)
	}
	channeld.StartProfiling()
	channeld.InitLogs()
	channeld.InitMetrics()
	channeld.InitConnections(channeld.GlobalSettings.ServerFSM, channeld.GlobalSettings.ClientFSM)
	channeld.InitChannels()

	channeld.InitSpatialController()

	unreal.InitMessageHandlers()

	// Setup Prometheus
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8080", nil)

	go channeld.StartListening(channeldpb.ConnectionType_SERVER, channeld.GlobalSettings.ServerNetwork, channeld.GlobalSettings.ServerAddress)

	// After the Master server owned the GLOBAL channel, the client connection should be listened.*/
	<-channeld.Event_GlobalChannelPossessed.Wait()
	channeld.StartListening(channeldpb.ConnectionType_CLIENT, channeld.GlobalSettings.ClientNetwork, channeld.GlobalSettings.ClientAddress)
}
