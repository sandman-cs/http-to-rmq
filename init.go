package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

// Configuration File Opjects
type configuration struct {
	ServerName           string
	AppName              string
	AppVer               string
	SrvPort              string
	Broker               string
	BrokerUser           string
	BrokerPwd            string
	BrokerExchange       string
	BrokerVhost          string
	ChannelSize          int
	ChannelCount         int
	ProxyHTTPSListenPort string //Listening port ofr SSL traffic, default port is 443
	LogLevel             string
	CrtFile              string //Path to cert file for SSL
	KeyFile              string //Path to key file for SSL
	AuthToken            []string
}

var (
	conf configuration
)

func init() {

	//Load Default Configuration Values
	conf.AppName = "Go - http-to-rmq"
	conf.AppVer = "1.0"
	conf.ServerName, _ = os.Hostname()
	conf.ChannelSize = 2048
	conf.SrvPort = "82"
	conf.LogLevel = "info"
	conf.Broker = "127.0.0.1"
	conf.BrokerUser = "guest"
	conf.BrokerPwd = "guest"
	conf.BrokerExchange = "amq.topic"
	conf.BrokerVhost = "/"
	conf.ChannelCount = 4
	conf.AuthToken = append(conf.AuthToken, "ItFDqpKDuEuJGS27+2m5bQ==")
	conf.ProxyHTTPSListenPort = "443"

	//Load Configuration Data
	dat, err := ioutil.ReadFile("conf.json")
	checkError(err)
	err = json.Unmarshal(dat, &conf)
	checkError(err)

	//fmt.Println("Config: ", conf)

	messages = make(chan chanToRabbit, conf.ChannelSize)

	// create the rabbitmq error channel
	rabbitCloseError = make(chan *amqp.Error)

	// run the callback in a separate thread
	go rabbitConnector(fmt.Sprint("amqp://" + conf.BrokerUser + ":" + conf.BrokerPwd + "@" + conf.Broker + conf.BrokerVhost))

	// establish the rabbitmq connection by sending
	// an error and thus calling the error callback
	rabbitCloseError <- amqp.ErrClosed

	for rabbitConn == nil {
		log.Println("Waiting for connection to rabbitmq...")
		time.Sleep(time.Second * 1)
	}

	for i := 0; i < conf.ChannelCount; i++ {
		go func() {
			for {
				chanPubToRabbit()
				time.Sleep(time.Second * 5)
			}
		}()
	}

}

func checkError(err error) {
	if err != nil {
		log.Println("Error: ", err)
	}
}
