package main

import (
	"flag"
	"fmt"
	"time"
	"sync"

	"github.com/golang/glog"
	"github.com/julienschmidt/httprouter"
	"net/http"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var configDir string
var configBind string
var _client MQTT.Client

var dataMap map[string]string
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

	dataMap = make(map[string]string)

	clientStop := make(chan struct{})
	defer close(clientStop)
	go mqttConnect(clientStop)

	router := httprouter.New()

	router.GET("/*queue", queueGet)
	router.POST("/*queue", queuePost)
	//router.PUT("/*queue", queuePost)

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
	var response string
	topic := params.ByName("queue")

	mux.RLock()
	if len(topic) > 1{
		response = string(dataMap[topic])
	}else{
		response = fmt.Sprintf("%#v\n", dataMap)
	}
	mux.RUnlock()

	w.WriteHeader(200)
	fmt.Fprint(w, string(response))
}

func queuePost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
        topic := params.ByName("queue")

	if len(topic) <= 1{
		w.WriteHeader(401)
		fmt.Fprint(w, "denied")
		return
	}

        message := "from gohome: " + time.Now().String()
        _client.Publish(topic, 0, false, message)

        w.WriteHeader(200)
        fmt.Fprint(w, "ok")
}

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
	mux.RLock()
	dataMap[message.Topic()] = string(message.Payload())
	mux.RUnlock()
}

