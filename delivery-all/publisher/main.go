package main

import (
	"encoding/json"
	"log"
	"nats-go/model"
	"strconv"

	"github.com/nats-io/nats.go"
)

const (
	subject1 = "subject1"
)

func main() {

	// Connect to NATS
	opt, err := nats.NkeyOptionFromSeed("../../nkey-cert.yaml")
	if err != nil {
		log.Fatalln("error nkey: ", err)
	}
	nc, err := nats.Connect("nats://localhost:4223", opt)
	if err != nil {
		log.Fatalln(err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalln(err)
	}

	err = createOrder(js)
	if err != nil {
		log.Fatalln(err)
	}
}

func createOrder(js nats.JetStreamContext) error {
	var order model.Order
	for i := 1; i <= 2; i++ {
		order = model.Order{
			OrderID:    i,
			CustomerID: "Delivery all-" + strconv.Itoa(i),
			Status:     "created",
		}
		orderJSON, _ := json.Marshal(order)
		_, err := js.Publish(subject1, orderJSON)
		if err != nil {
			return err
		}
		log.Printf("Order with OrderID:%d has been published\n", i)
	}
	return nil
}
