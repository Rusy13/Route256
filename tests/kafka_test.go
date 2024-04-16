package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"testing"
	"time"

	"Homework/internal/app/answer"
	"Homework/internal/app/payment"
	"Homework/internal/infrastructure/kafka"
)

func TestKafkaProducer(t *testing.T) {
	var testbroker = []string{

		"127.0.0.1:9093",
	}
	producerExample(testbroker, "TestURL", "TestMethod")
}

func producerExample(brokers []string, URL string, Method string) {

	kafkaProducer, err := kafka.NewProducer(brokers)
	if err != nil {
		fmt.Println(err)
	}

	answerService := answer.NewService(
		answer.NewRepository("some connector"),
		answer.NewKafkaSender(kafkaProducer, "testtopic"),
	)

	answerService.Verify(URL, Method, true, true)

	err = kafkaProducer.Close()
	if err != nil {
		fmt.Println("Close producers error ", err)
	}
}

func TestKafkaCons(t *testing.T) {
	var testbroker = []string{
		"127.0.0.1:9093",
	}

	outputChannel := make(chan string)
	go func() {
		ConsumerExample(testbroker, outputChannel)
	}()

	// Ожидание для установки соединения консьюмером
	time.Sleep(time.Second * 5)
	producerExample(testbroker, "TestURL", "TestMethod")
	time.Sleep(time.Second * 5)
	// Читаем вывод из канала
	output := <-outputChannel
	fmt.Println(output)
	ans := "Received Key: TestURL Value: {AnswerURL:TestURL Method: Success:true}"
	if ans != output {
		t.Errorf("Ошибка соответствия")
	}

}

type paymentMessage struct {
	AnswerURL string
	Method    string
	Success   bool
}

func ConsumerExample(brokers []string, outputChannel chan string) {
	kafkaConsumer, err := kafka.NewConsumer(brokers)
	if err != nil {
		fmt.Println(err)
	}

	// обработчики по каждому из топиков
	handlers := map[string]payment.HandleFunc{
		"testtopic": func(message *sarama.ConsumerMessage) {
			pm := paymentMessage{}
			err = json.Unmarshal(message.Value, &pm)
			if err != nil {
				fmt.Println("Consumer error", err)
			}

			var output string
			// Записываем вывод в переменную
			output = fmt.Sprintf("Received Key: %s Value: %+v", string(message.Key), pm)
			outputChannel <- output

		},
	}

	payments := payment.NewService(
		payment.NewReceiver(kafkaConsumer, handlers),
	)

	// При условии одного инстанса подходит идеально
	payments.StartConsume("testtopic")

	<-context.TODO().Done()
}
