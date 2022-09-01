package repository

import (
	"line/data"
	"line/models"
)

type MessageRepo struct {
	db *GenericRepository[data.MessageRecordPo]
}

/*
建立訊息資料存取層
*/
func NewMessageRepo() *MessageRepo {

	db := NewGenericRepository[data.MessageRecordPo]("messageRecord")
	return &MessageRepo{
		db,
	}

}

/*新增訊息*/
func (repo MessageRepo) Insert(model *data.MessageRecordPo) {

	repo.db.Insert(model)

}

/*
取得訊息列表
*/
func (repo MessageRepo) GetList(pageInfo *models.PageModel, filter interface{}) (*models.PageResult, *[]data.MessageRecordPo) {

	return repo.db.GetListPage(pageInfo, filter)
}
