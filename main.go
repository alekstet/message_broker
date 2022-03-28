package main

import (
	"flag"
	"net/http"

	"github.com/alekstet/message_broker/conf"
	"github.com/alekstet/message_broker/endpoint"
)

func main() {
	topics := conf.ReadConfig()

	for _, topic := range topics {
		d := endpoint.New(topic)
		http.HandleFunc(d.Url, d.Endpoint)
	}

	ip := flag.String("port", ":8080", "port")
	flag.Parse()
	http.ListenAndServe(*ip, nil)
}
