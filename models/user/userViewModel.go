package user

import "line/data"

type UserViewModel struct {
	UserID     string
	UserName   string
	UserPicUrl string
}

func (model *UserViewModel) MapToViewModel(data *data.UserPo) {
	model.UserID = data.UserID
	model.UserName = data.UserName
	model.UserPicUrl = data.UserPicUrl

}
