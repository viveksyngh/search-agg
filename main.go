package main

import "log"

func main() {

	server, err := CreateNewServer()
	if err != nil {
		log.Fatal(err.Error())
	}

	server.routes()
	server.Start()
}
