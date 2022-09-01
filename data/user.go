package data

type UserPo struct {
	UserID     string `bson:"_id,omitempty"`
	UserName   string `bson:"userName,omitempty"`
	UserPicUrl string `bson:"userPicUrl,omitempty"`
}

func NewUser(userId, userName, userPicUrl string) *UserPo {
	return &UserPo{
		UserID:     userId,
		UserName:   userName,
		UserPicUrl: userPicUrl,
	}
}
