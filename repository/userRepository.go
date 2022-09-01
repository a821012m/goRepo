package repository

import (
	"line/data"
	"line/models"

	"go.mongodb.org/mongo-driver/bson"
)

type UserRepo struct {
	db *GenericRepository[data.UserPo]
}

/*
建立使用者資料存取層
*/
func NewUserRepo() *UserRepo {

	db := NewGenericRepository[data.UserPo]("user")
	return &UserRepo{
		db,
	}

}

/* 取得使用者列表*/
func (repo *UserRepo) GetList(pageInfo *models.PageModel, name string) (*models.PageResult, *[]data.UserPo) {
	query := bson.M{}
	if len(name) > 0 {
		query["userName"] = bson.M{"$regex": name}
	}

	return repo.db.GetListPage(pageInfo, query)
}

/*
新增使用者
*/
func (repo *UserRepo) Insert(model *data.UserPo) {

	existsUser := repo.db.Get(&bson.D{{Key: "_id", Value: model.UserID}})

	if existsUser == nil {
		repo.db.Insert(model)
	}
}

/*取得單一使用者*/
func (repo *UserRepo) Get(userId string) *data.UserPo {
	existsUser := repo.db.Get(&bson.D{{Key: "_id", Value: userId}})
	return existsUser
}
