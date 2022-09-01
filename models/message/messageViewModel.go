package message

type MessageViewModel struct {
	UserId     string
	UserName   string
	UserPicUrl string
	Text       string
}

func NewMessageViewModel(UserId string,
	UserName string,
	UserPicUrl string,
	Text string) *MessageViewModel {

	return &MessageViewModel{
		UserId:     UserId,
		UserName:   UserName,
		UserPicUrl: UserPicUrl,
		Text:       Text,
	}
}
