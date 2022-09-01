package user

import "line/models"

type UserListDto struct {
	UserName string `uri:"name" `
	models.PageModel
}
