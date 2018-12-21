package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/taviti/caldav-go/caldav"
)


func main() {
	// create a reference to your CalDAV-compliant server
	server, err := caldav.NewServer("URL")

	if err != nil {
		log.Fatal(err)
	}

	// create a CalDAV client to speak to the server
	var client = caldav.NewClient(server, http.DefaultClient)

	// start executing requests!
	err = client.ValidateServer("")
	log.Fatal(err)
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", rootHandler)

	panic(http.ListenAndServe(":8080", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		fmt.Println("Could not open file.", err)
	}
	fmt.Fprintf(w, "%s", content)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	var parsed []map[string]interface{}

	content, err := ioutil.ReadFile("monitor.json")
	if err != nil {
		fmt.Println("Could not open file.", err)
	}

	err = json.Unmarshal(content, &parsed)
	if err != nil {
		fmt.Println("Could not open file.", err)
	}

	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}

	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	conn.WriteJSON(parsed)

	go echo(conn)
}

func echo(conn *websocket.Conn) {
	for {
		var m []map[string]interface{}

		err := conn.ReadJSON(&m)
		if err != nil {
			fmt.Println("Error reading json.", err)
		}

		fmt.Printf("Got message: %#v\n", m)

		if err = conn.WriteJSON(m); err != nil {
			fmt.Println(err)
		}
	}
}
