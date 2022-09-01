package data

import "go.mongodb.org/mongo-driver/bson/primitive"

type MessageRecordPo struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	UserID          string             `bson:"userID,omitempty"`
	Text            string             `bson:"text,omitempty"`
	FullMessageJson string             `bson:"fullMessageJson,omitempty"`
}

func NewMessage(userId, text, fullMessageJson string) *MessageRecordPo {
	return &MessageRecordPo{
		UserID:          userId,
		Text:            text,
		FullMessageJson: fullMessageJson,
	}
}
