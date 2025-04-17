package mq

import (
	"eshop_main/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	Conn *amqp.Connection
	url  = "amqp://admin:admin@117.72.72.114:5672/"
)

func init() {
	var err error
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	Conn = conn
	ch, err := conn.Channel()
	defer ch.Close()
	// 延迟交换机
	err = createExchange(ch, "exchange_eshop_main", "direct")
	err = createQueue(ch, "update_surplus_queue")
	err = ch.QueueBind(
		"update_surplus_queue", // queue name
		"update_surplus",       // routing key
		"exchange_eshop_main",  // exchange name
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		log.Errorf("err: %v", err)
	}
	// 消费者
	go ConsumeUpdateSurplusMessage()
	return
}

func createQueue(ch *amqp.Channel, name string) (err error) {
	_, err = ch.QueueDeclare(
		name,  // queue name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Errorf("err: %v", err)
	}
	return
}

func createExchange(ch *amqp.Channel, name string, kind string) (err error) {
	err = ch.ExchangeDeclare(
		name,  // exchange name
		kind,  // exchange type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		amqp.Table{
			"x-delayed-type": "direct",
		},
	)
	if err != nil {
		log.Errorf("err: %v", err)
	}
	return
}
