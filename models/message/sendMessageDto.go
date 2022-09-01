package message

type SendMessageDto struct {
	UserId string `form:"UserId" json:"UserId" binding:"required"`
	Text   string `form:"Text" json:"Text" binding:"required"`
}
