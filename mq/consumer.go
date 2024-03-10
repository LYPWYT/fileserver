package mq

import "log"

var done chan bool

// StartConsume: 开始监听队列，获取消息
func StartConsume(qName, cName string, callback func(msg []byte) bool) {
	// 1.通过channel.Consume获得消息信道
	msgs, err := channel.Consume(
		qName,
		cName,
		true,  // 自动回复确认收到的信息
		false, // 是否是唯一消费者这里选择不唯一
		false,
		false, // 客户端在创建请求过后会等待rabbitmq返回一些信息
		nil,
	)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// 2.循环获取队列的消息

	done = make(chan bool)

	go func() {
		for msg := range msgs {
			// 3.调用callback方法来处理新的消息
			processSuc := callback(msg.Body)
			if !processSuc {
				// TODO: 将任务写到另一个队列，用于异常情况的重试
			}
		}
	}()

	// done没有新的消息进来，就会一直处于阻塞状态
	<-done

	// 关闭rabbitMQ通道
	channel.Close()
}

// StopConsume: 停止监听队列
func StopConsume() {
	done <- true
}
