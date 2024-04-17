package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

type paymentMessage struct {
	AnswerURL string
	Method    string
	Success   bool
}

type ConsumerGroup struct {
	ready chan bool
}

func NewConsumerGroup() ConsumerGroup {
	return ConsumerGroup{
		ready: make(chan bool),
	}
}

func (consumer *ConsumerGroup) Ready() <-chan bool {
	return consumer.ready
}

// Setup Начинаем новую сессию, до ConsumeClaim
func (consumer *ConsumerGroup) Setup(_ sarama.ConsumerGroupSession) error {
	close(consumer.ready)

	return nil
}

// Cleanup завершает сессию, после того, как все ConsumeClaim завершатся
func (consumer *ConsumerGroup) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim читаем до тех пор пока сессия не завершилась
func (consumer *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		log.Println("qwertu")
		select {
		case message := <-claim.Messages():

			pm := paymentMessage{}
			err := json.Unmarshal(message.Value, &pm)
			if err != nil {
				fmt.Println("Consumer group error", err)
			}

			log.Printf("Message claimed: URL = %v,Message claimed: Method = %v, timestamp = %v, topic = %s",
				pm.AnswerURL,
				pm.Method,
				message.Timestamp,
				message.Topic,
			)

			// коммит сообщения "руками"
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
