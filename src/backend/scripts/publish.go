package main

import (
	"log"
	"os"

	stan "github.com/nats-io/stan.go"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: publish <json-file>")
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	sc, err := stan.Connect("test-cluster", "publisher", stan.NatsURL("nats://nats:4222"))
	if err != nil {
		log.Fatal("NATS connect error:", err)
	}
	defer sc.Close()

	if err := sc.Publish("orders", data); err != nil {
		log.Fatal("publish error:", err)
	}
	log.Println("published")
}
