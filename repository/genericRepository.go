package repository

import (
	"context"
	"fmt"
	"line/data"
	"line/models"
	"log"
	"math"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type GenericRepository[T data.UserPo | data.MessageRecordPo] struct {
	client     *mongo.Client
	context    context.Context
	cancel     context.CancelFunc
	collection *mongo.Collection
}

func NewGenericRepository[T data.UserPo | data.MessageRecordPo](tableName string) *GenericRepository[T] {

	connectStr := viper.GetString("System.ConnectString")
	dbName := viper.GetString("System.DataBase")
	fmt.Println("db" + dbName)
	if len(connectStr) == 0 || len(dbName) == 0 {
		log.Panic("查無DB設定檔")
		panic("查無DB設定檔")

	}

	context, cancelFunc, client, collection := connect(connectStr, dbName, tableName)
	return &GenericRepository[T]{
		client, context, cancelFunc, collection,
	}
}

func connect(connectStr, dbName, tableName string) (context.Context, context.CancelFunc, *mongo.Client, *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	innerClient, err := mongo.Connect(ctx, options.Client().ApplyURI(connectStr))
	err = innerClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Panic(err.Error())
		panic(err)
	}
	collection := innerClient.Database(dbName).Collection(tableName)
	return ctx, cancel, innerClient, collection
}

func (repo *GenericRepository[T]) Get(filter *bson.D) *T {
	result := new(T)
	repo.context, repo.cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer repo.cancel()
	err := repo.collection.FindOne(repo.context, filter).Decode(result)

	if err == mongo.ErrNoDocuments {
		return nil
	} else {
		return result
	}
}

func (repo *GenericRepository[T]) GetListPage(pageInfo *models.PageModel, filter interface{}) (*models.PageResult, *[]T) {
	repo.context, repo.cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer repo.cancel()
	res := new(models.PageResult)

	result := make([]T, 0)

	count, err := repo.collection.CountDocuments(repo.context, filter)
	if err != nil {
		log.Panic(err.Error())
	}

	res.TotalCount = count
	res.TotalPage = int64(math.Ceil(float64(count) / float64(pageInfo.PageSize)))
	findOptions := options.Find()

	findOptions.SetSkip(pageInfo.PageSize * (pageInfo.Index - 1))
	findOptions.SetLimit(pageInfo.PageSize)

	cur, err := repo.collection.Find(repo.context, filter, findOptions)

	if err != nil {
		log.Panic(err.Error())
	}

	for cur.Next(repo.context) {
		tempData := new(T)
		err := cur.Decode(tempData)
		if err != nil {
			log.Panic(err.Error())
		} else {

			result = append(result, *tempData)

		}

	}

	if err := cur.Err(); err != nil {
		log.Panic(err.Error())
	}

	//Close the cursor once finished
	cur.Close(repo.context)
	res.PageInfo = *pageInfo

	return res, &result
}

func (repo *GenericRepository[T]) Insert(data interface{}) {
	repo.context, repo.cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer repo.cancel()
	insertResult, err := repo.collection.InsertOne(repo.context, data)

	if err != nil {
		log.Panic(err.Error())
	}

	fmt.Println(insertResult.InsertedID)

}
