package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"encoding/json"

	"github.com/golang/glog"
	"github.com/julienschmidt/httprouter"
	"net/http"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	"github.com/mbroome/gohome/pkg/persist"
)

type DataPoint struct {
	Value     string    `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

var configDir string
var configBind string
var _client MQTT.Client

var dataMap map[string]DataPoint
var mux sync.RWMutex

func main() {
	flag.StringVar(&configDir, "config", "", "Path to config dir")
	flag.StringVar(&configBind, "bind", "", "Interface:port to bind to")
	flag.Parse()

	if configDir == "" {
		configDir = "/etc/config/exporters"
	}
	if configBind == "" {
		configBind = "0.0.0.0:8080"
	}

	dataMap = make(map[string]DataPoint)

	clientStop := make(chan struct{})
	defer close(clientStop)
	go mqttConnect(clientStop)

	router := httprouter.New()

	router.GET("/topic/*queue", queueGet)
	router.PUT("/topic/*queue", queuePut)
	router.GET("/list/*queue", queueList)

	glog.Fatal(http.ListenAndServe(configBind, router))

}

func mqttConnect(c chan struct{}) {
	qos := 0
	server := "tcp://127.0.0.1:1883"
	clientid := "gohome"
	topic := "#"

	connOpts := MQTT.NewClientOptions().AddBroker(server).SetClientID(clientid).SetCleanSession(true)

	connOpts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(topic, byte(qos), onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	_client = MQTT.NewClient(connOpts)
	if token := _client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return
	}

	fmt.Printf("Connected to %s\n", server)

	<-c
}

func queueGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var response []byte
	topic := params.ByName("queue")

	mux.RLock()
	if len(topic) > 1 {
		response, _ = json.Marshal(dataMap[topic])
	} else {
		response, _ = json.Marshal(dataMap)
	}
	mux.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}

func queuePut(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	topic := params.ByName("queue")

	data, _ := ioutil.ReadAll(r.Body)

	if len(topic) <= 1 {
		w.WriteHeader(401)
		fmt.Fprint(w, "denied")
		return
	}

	_client.Publish(topic, 0, false, string(data))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprint(w, "{\"status\":\"ok\"}")
}

func queueList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	keys := []string{}
	mux.RLock()
	for k := range dataMap {
		keys = append(keys, k)
	}
	mux.RUnlock()

	out, err := json.Marshal(keys)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(out)
}

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	var rec DataPoint
	rec.Value = string(message.Payload())
	rec.Timestamp = time.Now()

	mux.RLock()
	dataMap[message.Topic()] = rec
	if err := persist.Save("./file.tmp", dataMap); err != nil {
		fmt.Print(err)
	}

	mux.RUnlock()
}
