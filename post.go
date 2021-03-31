package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func postHandlerWithToken(w http.ResponseWriter, r *http.Request) {
	var msg chanToRabbit
	msg.route = "p.mos.pmg.usage"
	authToken := r.Header.Get("X-Auth-Token")
	if stringContains(authToken, conf.AuthToken) {
		log.Println("Recieved the Auth Token")
	} else {
		log.Println("Didn't recieve the Auth Token")
		r.Body.Close()
		return
	}
	receivedJSON, err := ioutil.ReadAll(r.Body) //This reads raw request body
	if err != nil {
		io.WriteString(w, "Invalid...\n")
		msg.route = "p.mos.pmg.usage.bad"
		msg.payload = string(receivedJSON)
		messages <- msg
	} else {

		msg.payload = string(receivedJSON)
		messages <- msg
	}
	io.WriteString(w, "ok\n")

	r.Body.Close()
}
