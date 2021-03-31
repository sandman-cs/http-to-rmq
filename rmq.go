package main

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

type chanToRabbit struct {
	payload string
	route   string
}

var (
	pubSuccess       int
	pubError         int
	messages         chan chanToRabbit
	rabbitConn       *amqp.Connection
	rabbitCloseError chan *amqp.Error
)

// Try to connect to the RabbitMQ server as
// long as it takes to establish a connection
//
func connectToRabbitMQ(uri string) *amqp.Connection {
	for {
		conn, err := amqp.Dial(uri)

		if err == nil {
			return conn
		}

		checkError(err)
		log.Println("Trying to reconnect to RabbitMQ")
		time.Sleep(500 * time.Millisecond)
	}
}

// re-establish the connection to RabbitMQ in case
// the connection has died
//
func rabbitConnector(uri string) {
	var rabbitErr *amqp.Error

	for {
		rabbitErr = <-rabbitCloseError
		if rabbitErr != nil {
			log.Println("Connecting to RabbitMQ:")
			rabbitConn = connectToRabbitMQ(uri)
			rabbitCloseError = make(chan *amqp.Error)
			rabbitConn.NotifyClose(rabbitCloseError)
		}
	}
}

func chanPubToRabbit() {

	ch, err := rabbitConn.Channel()
	checkError(err)
	if err == nil {
		defer ch.Close()
		for {
			msg := <-messages
			//fmt.Println("Exchange: ", conf.BrokerExchange, " Route: ", msg.route, " msg: ", msg.payload)
			err = ch.Publish(
				conf.BrokerExchange, // exchange
				msg.route,           // routing key
				false,               // mandatory
				false,               // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(msg.payload),
				})
			checkError(err)
			if err != nil {
				messages <- msg
				pubError++
				break
			} else {
				pubSuccess++
			}
		}
	}
}
