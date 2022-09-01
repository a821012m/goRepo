package service

import (
	"encoding/json"
	"line/data"
	"line/models"
	"line/models/message"
	"line/models/user"
	"line/repository"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
)

var channelSecret string
var accessToken string

type LineService struct {
	lineClient  *linebot.Client
	userRepo    *repository.UserRepo
	messageRepo *repository.MessageRepo
}

/*
建立物件
*/
func NewLineService() *LineService {

	channelSecret = viper.GetString("Line.ChannelSecret")
	accessToken = viper.GetString("Line.AccessToken")
	if len(channelSecret) == 0 || len(accessToken) == 0 {
		panic("line服務 AccessToken 或 ChannelSecret 未設定")
	}

	return &LineService{
		lineServiceInit(),
		repository.NewUserRepo(),
		repository.NewMessageRepo(),
	}
}

func lineServiceInit() *linebot.Client {
	client := &http.Client{}
	bot, err := linebot.New(channelSecret, accessToken, linebot.WithHTTPClient(client))
	if err != nil {
		log.Panic(err.Error())
		panic(err)
	} else {
		return bot
	}

}

/* 寄出訊息給使用者*/
func (service *LineService) SendMessage(userId, text string) bool {
	msg := linebot.NewTextMessage(text)

	var messages []linebot.SendingMessage
	messages = append(messages, msg)

	_, err := service.lineClient.PushMessage(userId, messages...).Do()
	if err != nil {
		log.Panic(err.Error())
	}
	return err == nil

	// fmt.Printf("temp.RequestID: %v\n", temp.RequestID)
}

/*
接收line 訊息並儲存
*/
func (service *LineService) ReceiveText(w http.ResponseWriter, r *http.Request) {

	events, err := service.lineClient.ParseRequest(r)
	if err != nil {
		log.Panic(err.Error())
	}
	for _, event := range events {
		userProfileRes := service.GetUserProfileFromLine(event.Source.UserID)
		service.saveUser(userProfileRes)

		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				service.saveTextMessage(message, event.Source.UserID)

				if _, err = service.lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("已儲存文字訊息")).Do(); err != nil {
					log.Panic(err.Error())
				}
			case *linebot.StickerMessage:
				service.saveStickerMessage(message, event.Source.UserID)

				if _, err = service.lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("已儲存貼圖")).Do(); err != nil {
					log.Panic(err.Error())
				}
			default:
				out, err := json.Marshal(message)
				if err != nil {
					log.Panic(err.Error())
				}
				service.saveMessage(event.Source.UserID, string(event.Type), string(out))
				if _, err = service.lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("尚未實作此訊息類型")).Do(); err != nil {
					log.Panic(err.Error())
				}
			}
		}
	}

}

func (service *LineService) GetMessageRecords(pageInfo *models.PageModel, userId string) *models.PageResult {
	query := bson.M{}
	if len(userId) > 0 {
		query["userID"] = bson.M{"$eq": userId}
	}

	pageResult, datas := service.messageRepo.GetList(pageInfo, query)
	viewModels := make([]message.MessageViewModel, 0)
	for _, item := range *datas {
		userData := service.userRepo.Get(item.UserID)
		viewModels = append(viewModels, *message.NewMessageViewModel(userData.UserID, userData.UserName, userData.UserPicUrl, item.Text))
	}
	pageResult.Datas = viewModels
	return pageResult
}

/*
取得line使用者資料
*/
func (service *LineService) GetUserProfileFromLine(userId string) *linebot.UserProfileResponse {
	userRes, err := service.lineClient.GetProfile(userId).Do()

	if err != nil {
		log.Panic(err.Error())
		// panic(err)
	}
	return userRes
}

func (service *LineService) GetUserList(pageInfo *models.PageModel, name string) *models.PageResult {
	pageRes, datas := service.userRepo.GetList(pageInfo, name)
	userViewLs := make([]user.UserViewModel, len(*datas))
	for index := 0; index < len(userViewLs); index++ {
		newViewModel := new(user.UserViewModel)
		item := (*datas)[index]
		newViewModel.MapToViewModel(&item)
		userViewLs[index] = *newViewModel
	}

	pageRes.Datas = userViewLs
	return pageRes
}

/*
儲存使用者資料
*/
func (service *LineService) saveUser(lineUser *linebot.UserProfileResponse) {

	user := data.NewUser(lineUser.UserID, lineUser.DisplayName, lineUser.PictureURL)
	service.userRepo.Insert(user)
}

/*
儲存文字訊息
*/
func (service *LineService) saveTextMessage(message *linebot.TextMessage, userId string) {
	out, err := json.Marshal(message)
	if err != nil {
		log.Panic(err.Error())
	}
	service.saveMessage(userId, message.Text, string(out))
}

/*
儲存貼圖訊息
*/
func (service *LineService) saveStickerMessage(message *linebot.StickerMessage, userId string) {
	out, err := json.Marshal(message)
	if err != nil {
		log.Panic(err.Error())
	}
	service.saveMessage(userId, "貼圖id:"+message.StickerID, string(out))
}

/*
儲存訊息到資料庫
*/
func (service *LineService) saveMessage(userId, text, fullMessageJson string) {

	newMessageRecord := data.NewMessage(userId, text, fullMessageJson)

	service.messageRepo.Insert(newMessageRecord)
}
