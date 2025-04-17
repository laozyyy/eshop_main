package mq

import (
	"context"
	"encoding/json"
	"eshop_main/cache"
	"eshop_main/database"
	"eshop_main/log"
	"fmt"
)

func ConsumeUpdateSurplusMessage() {
	ch, err := Conn.Channel()
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	defer ch.Close()
	// 消费消息

	msgs, err := ch.Consume(
		"update_surplus_queue",          // 队列名称
		"update_surplus_queue_consumer", // 消费者标签
		false,                           // 是否自动确认消息
		false,                           // 是否独占消费者（仅限于本连接）
		false,                           // 是否阻塞等待服务器确认
		false,                           // 是否使用内部排他队列
		nil,                             // 其他参数
	)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	log.Infof("update_surplus_queue_consumer消费者启动")
	for msg := range msgs {
		log.Infof("update_surplus_queue_consumer消费者收到消息: %s", string(msg.Body))
		// 更新库
		message := make(map[string]interface{})
		err = json.Unmarshal(msg.Body, &message)
		if err != nil {
			log.Errorf("err: %v", err)
			return
		}
		sku := message["sku"].(string)
		aId := message["activity_id"].(string)
		surplusKey := fmt.Sprintf("seckill_surplus:%s", aId)
		surplus, err := cache.Client.Get(context.Background(), surplusKey).Result()
		if err != nil {
			log.Errorf("err: %v", err)
			continue
		}
		err = database.UpdateSurplus(nil, sku, surplus)
		if err != nil {
			log.Errorf("err: %v", err)
			continue
		}
		// 手动确认消息已被消费
		msg.Ack(false)
	}
}

func SendLotterySurplusDecr(id string) {

}
