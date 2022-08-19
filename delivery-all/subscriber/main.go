package main

import (
	"encoding/json"
	"fmt"
	"log"
	"nats-go/model"
	"runtime"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	stream1  = "stream1"
	durable1 = "durable1"
	queue1   = "queue1"
	subject1 = "subject1"
	deliver1 = "deliver1"
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

	t := time.Now()
	fmt.Println(t)
	createStream(js)
	createConsumer(js, t)
	createSub(js, t)

	runtime.Goexit()
}

// create stream
func createStream(js nats.JetStreamContext) {
	info, _ := js.StreamInfo(stream1)
	if info != nil {
		js.DeleteStream(stream1)
	}

	cfg := nats.StreamConfig{
		Name:      stream1,
		Retention: nats.InterestPolicy,
		Replicas:  1,
		Subjects:  []string{subject1},
	}
	_, err := js.AddStream(&cfg)
	if err != nil {
		log.Fatal("create stream error: ", err.Error())
	}

	log.Println("=========================== create stream success")
}

// create consumer
func createConsumer(js nats.JetStreamContext, now time.Time) {
	info, _ := js.ConsumerInfo(stream1, durable1)
	if info != nil {
		js.DeleteConsumer(stream1, durable1)
	}
	// timeNow := time.Now()
	cfg := nats.ConsumerConfig{
		Durable:        durable1,
		FilterSubject:  subject1,
		MaxDeliver:     3,
		DeliverSubject: deliver1,
		AckWait:        time.Second * 5,
		DeliverPolicy:  nats.DeliverAllPolicy,
		DeliverGroup:   queue1,
		AckPolicy:      nats.AckExplicitPolicy,
	}
	_, err := js.AddConsumer(stream1, &cfg)
	if err != nil {
		log.Fatal("create consumer error: ", err.Error())
	}
	log.Println("=========================== create consumer success")

}

// register queue-subscribe
func createSub(js nats.JetStreamContext, now time.Time) {
	_, err := js.QueueSubscribe(subject1, queue1, func(msg *nats.Msg) {
		var order model.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("monitor service subscribes from subject:%s\n", msg.Subject)
		log.Printf("OrderID:%d, CustomerID: %s, Status:%s\n", order.OrderID, order.CustomerID, order.Status)
		fmt.Println("---------------")
		// msg.Ack()
	}, nats.Bind(stream1, durable1),
		nats.ManualAck(),
		nats.MaxDeliver(3),
		nats.DeliverSubject(deliver1),
		nats.AckWait(time.Second*5),
	)

	if err != nil {
		log.Fatal("create createSub error: ", err.Error())
	}

	log.Println("=========================== create createSub success")

}
