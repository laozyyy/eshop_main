package mq

import (
	"encoding/json"
	"eshop_main/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SendMessage(exchange string, routingKey string, publish amqp.Publishing) error {
	ch, err := Conn.Channel()
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	defer ch.Close()
	err = ch.Publish(
		exchange,   // exchange name
		routingKey, // routing key
		false,
		false,
		publish,
	)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	log.Infof("发送消息成功: %v", string(publish.Body))
	return nil
}

func SendSurplusDecr(sku string, aId string) {
	message := map[string]interface{}{
		"sku":         sku,
		"activity_id": aId,
	}
	marshal, _ := json.Marshal(message)
	publish := amqp.Publishing{
		ContentType: "text/plain",
		Body:        marshal,
	}
	log.Infof("发送购物车写库消息")
	_ = SendMessage("exchange_eshop_main", "update_surplus", publish)
}
