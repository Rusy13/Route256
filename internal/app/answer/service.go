package answer

import "fmt"

type Sender interface {
	sendAsyncMessage(message paymentMessage) error
	sendMessage(message paymentMessage) error
	sendMessages(messages []paymentMessage) error
}

type Repository interface {
	GetAnswer(URL string, Method string) *Answer
	VerifyAnswer(URL string, Method string, success bool) error
}

type Service struct {
	repo   Repository
	sender Sender
}

func NewService(repo Repository, sender Sender) *Service {
	return &Service{
		repo:   repo,
		sender: sender,
	}
}

func (s Service) Verify(answerURL string, answerMethod string, success bool, sync bool) {
	answer := s.repo.GetAnswer(answerURL, answerMethod)
	err := s.repo.VerifyAnswer(answerURL, answerMethod, success)
	if err != nil {
		fmt.Println("Error in verify", err)
	}

	if sync {
		err = s.sender.sendMessage(
			paymentMessage{
				answer.URL,
				answer.Method,
				success,
			},
		)

		if err != nil {
			fmt.Println("Send sync message error: ", err)
		}

		return
	}

	//err = s.sender.sendAsyncMessage(
	//	paymentMessage{
	//		answer.ID,
	//		answer.userID,
	//		answer.sum,
	//		success,
	//	},
	//)
	//
	//if err != nil {
	//	fmt.Println("Send async message error: ", err)
	//}
}

//func (s Service) VerifyBatch(answerIDs []int, success bool) {
//	var messages []paymentMessage
//
//	for _, ID := range answerIDs {
//		answer := s.repo.GetAnswer(ID)
//		err := s.repo.VerifyAnswer(ID, success)
//		if err != nil {
//			fmt.Println("Error in verify", err)
//		}
//
//		messages = append(messages, paymentMessage{
//			answer.ID,
//			answer.userID,
//			answer.sum,
//			success,
//		})
//	}
//
//	err := s.sender.sendMessages(messages)
//
//	if err != nil {
//		fmt.Println("Send message error: ", err)
//	}
//}
