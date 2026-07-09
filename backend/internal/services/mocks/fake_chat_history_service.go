package mocks

type FakeChatHistoryService struct {
	SavedQuestion string
	SavedAnswer   string
	Err           error
}

func (s *FakeChatHistoryService) Save(question, answer string) error {
	s.SavedQuestion = question
	s.SavedAnswer = answer

	return s.Err
}
