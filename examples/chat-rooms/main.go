package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"channeld.clewcat.com/channeld/pkg/channeld"
	"channeld.clewcat.com/channeld/proto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// The addresses should only include the port for running in the docker containers.
var webAddr = flag.String("web", ":8080", "http service address")
var wsAddr = flag.String("ws", ":12108", "websocket service address")

func handleChanneldProto(w http.ResponseWriter, r *http.Request) {
	bytes, err := os.ReadFile("../../proto/channeld.proto")
	if err != nil {
		log.Panic(err)
	}
	w.Write(bytes)
}

func handleChatProto(w http.ResponseWriter, r *http.Request) {
	bytes, err := os.ReadFile("../../proto/example_chat_rooms.proto")
	if err != nil {
		log.Panic(err)
	}
	w.Write(bytes)
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	homeTemplate, err := template.ParseFiles("./main.html")
	if err != nil {
		log.Panic(err)
	}

	localAddr := r.Host[:strings.Index(r.Host, ":")]
	//log.Println(localAddr)
	homeTemplate.Execute(w, localAddr+*wsAddr)
}

func main() {

	http.HandleFunc("/", handleMain)
	http.HandleFunc("/proto", handleChanneldProto)
	http.HandleFunc("/proto/chat", handleChatProto)
	http.Handle("/metrics", promhttp.Handler())

	channeld.InitLogsAndMetrics()
	channeld.InitConnections("../../config/server_authoratative_fsm.json", "../../config/client_authoratative_fsm.json")
	channeld.InitChannels()
	channeld.GetChannel(channeld.GlobalChannelId).InitData(
		&proto.ChatChannelData{ChatMessages: []*proto.ChatMessage{
			{Sender: "System", SendTime: time.Now().Unix(), Content: "Welcome!"},
		}},
		&channeld.DataMergeOptions{ListSizeLimit: 100},
	)
	//channeld.SetWebSocketTrustedOrigins(["localhost"])
	go channeld.StartListening(channeld.CLIENT, "ws", *wsAddr)

	log.Fatal(http.ListenAndServe(*webAddr, nil))

}
