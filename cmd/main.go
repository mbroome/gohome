package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
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
	ID        string    `json:"id"`
	Value     string    `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type DeviceResponseRecord struct {
	Name      string    `json:"name"`
	ID        string    `json:"id"`
	Value     string    `json:"value"`
	Group     string    `json:"group"`
	Timestamp time.Time `json:"timestamp"`
}

type Config struct {
	Devices []DeviceDetails `json:"devices"`
}

type DeviceDetails struct {
	Name  string `json:"name"`
	Write string `json:"write"`
	Read  string `json:"read"`
	Group string `json:"group"`
}

var configFile string
var configBind string
var _client MQTT.Client
var deviceConfig Config

var dataMap map[string]DataPoint
var mux sync.RWMutex

func main() {
	flag.StringVar(&configFile, "config", "", "Path to config file")
	flag.StringVar(&configBind, "bind", "", "Interface:port to bind to")
	flag.Parse()

	if configFile == "" {
		configFile = "config.json"
	}
	if configBind == "" {
		configBind = "0.0.0.0:8080"
	}

	deviceConfig = LoadConfiguration(configFile)
	fmt.Printf("%#v\n", deviceConfig)

	dataMap = make(map[string]DataPoint)

	if err := persist.Load("./file.tmp", &dataMap); err != nil {
		fmt.Print(err)
	}

	clientStop := make(chan struct{})
	defer close(clientStop)
	go mqttConnect(clientStop)

	router := httprouter.New()

	router.GET("/data", queueGet)
	router.PUT("/command/*queue", queuePut)
	router.GET("/list/*queue", queueList)

	glog.Fatal(http.ListenAndServe(configBind, router))

}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func mqttConnect(c chan struct{}) {
	qos := 0
	//server := "tcp://127.0.0.1:1883"
	server := "tcp://192.168.1.79:1883"
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

func queueGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var response []byte
	var dataPoints []DeviceResponseRecord

	mux.RLock()
	for device := range deviceConfig.Devices {
		//fmt.Printf("##### device: %s\n", deviceConfig.Devices[device].Read)
		for point := range dataMap {
			//fmt.Printf("##### point: %s\n", dataMap[point].ID)
			if dataMap[point].ID == deviceConfig.Devices[device].Read {
				//fmt.Printf("############ found a matched read field: %s\n", dataMap[point].ID)
				var rr DeviceResponseRecord
				rr.Name = deviceConfig.Devices[device].Name
				rr.Group = deviceConfig.Devices[device].Group
				rr.ID = dataMap[point].ID
				rr.Value = dataMap[point].Value
				rr.Timestamp = dataMap[point].Timestamp
				dataPoints = append(dataPoints, rr)
			}
		}
	}
	mux.RUnlock()

	response, _ = json.Marshal(dataPoints)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	rec.ID = message.Topic()

	mux.RLock()

	for device := range deviceConfig.Devices {
		//fmt.Printf("##### device: %s\n", deviceConfig.Devices[device].Read)
		if (rec.ID == deviceConfig.Devices[device].Read) || (rec.ID == deviceConfig.Devices[device].Write) {

			if _, ok := dataMap[rec.ID]; ok {
				if dataMap[rec.ID].Value != rec.Value {
					fmt.Printf("Value changed: %s => %s\n", rec.ID, rec.Value)
				}
			} else {
				fmt.Printf("New value: %s => %s\n", rec.ID, rec.Value)
			}

			dataMap[rec.ID] = rec
			if err := persist.Save("./file.tmp", dataMap); err != nil {
				fmt.Print(err)
			}
		}
	}
	mux.RUnlock()
}
