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
	stream2  = "stream2"
	durable2 = "durable2"
	queue2   = "queue2"
	subject2 = "subject2"
	deliver2 = "deliver2"
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
	info, _ := js.StreamInfo(stream2)
	if info != nil {
		js.DeleteStream(stream2)
	}

	cfg := nats.StreamConfig{
		Name:      stream2,
		Retention: nats.InterestPolicy,
		Replicas:  1,
		Subjects:  []string{subject2},
	}
	_, err := js.AddStream(&cfg)
	if err != nil {
		log.Fatal("create stream error: ", err.Error())
	}

	log.Println("=========================== create stream success")
}

// create consumer
func createConsumer(js nats.JetStreamContext, now time.Time) {
	info, _ := js.ConsumerInfo(stream2, durable2)
	if info != nil {
		js.DeleteConsumer(stream2, durable2)
	}
	// timeNow := time.Now()
	cfg := nats.ConsumerConfig{
		Durable:        durable2,
		FilterSubject:  subject2,
		MaxDeliver:     3,
		DeliverSubject: deliver2,
		AckWait:        time.Second * 5,
		DeliverPolicy:  nats.DeliverNewPolicy,
		DeliverGroup:   queue2,
		AckPolicy:      nats.AckExplicitPolicy,
	}
	_, err := js.AddConsumer(stream2, &cfg)
	if err != nil {
		log.Fatal("create consumer error: ", err.Error())
	}
	log.Println("=========================== create consumer success")

}

// register queue-subscribe
func createSub(js nats.JetStreamContext, now time.Time) {
	_, err := js.QueueSubscribe(subject2, queue2, func(msg *nats.Msg) {
		var order model.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("monitor service subscribes from subject:%s\n", msg.Subject)
		log.Printf("OrderID:%d, CustomerID: %s, Status:%s\n", order.OrderID, order.CustomerID, order.Status)
		fmt.Println("---------------")
		// msg.Ack()
	}, nats.Bind(stream2, durable2),
		nats.ManualAck(),
		nats.MaxDeliver(3),
		nats.DeliverSubject(deliver2),
		nats.AckWait(time.Second*5),
	)

	if err != nil {
		log.Fatal("create createSub error: ", err.Error())
	}

	log.Println("=========================== create createSub success")

}
