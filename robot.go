package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
	yaml "gopkg.in/yaml.v2"

	"demo/utils"
)

// Config 定义了配置文件的结构
type Config struct {
	AppID uint64 `yaml:"appid"` //机器人的appid
	Token string `yaml:"token"` //机器人的token
}

var config Config
var api openapi.OpenAPI
var ctx context.Context
var channelId = "" //保存子频道的id

// 第一步： 获取机器人的配置信息，即机器人的appid和token
func init() {
	content, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Println("读取配置文件出错， err = ", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		log.Println("解析配置文件出错， err = ", err)
		os.Exit(1)
	}
	log.Println(config)
}
func main() {
	token := token.BotToken(config.AppID, config.Token) //生成token
	api = botgo.NewOpenAPI(token).WithTimeout(3 * time.Second)
	ctx = context.Background()
	ws, err := api.WS(ctx, nil, "") //websocket
	if err != nil {
		log.Fatalln("websocket错误， err = ", err)
		os.Exit(1)
	}

	var atMessage event.ATMessageEventHandler = atMessageEventHandler

	intent := websocket.RegisterHandlers(atMessage)     // 注册socket消息处理
	botgo.NewSessionManager().Start(ws, token, &intent) // 启动socket监听
}

// atMessageEventHandler 处理 @机器人 的消息
func atMessageEventHandler(event *dto.WSPayload, data *dto.WSATMessageData) error {
	channelId = data.ChannelID //当@机器人时，保存ChannelId，主动消息需要 channelId 才能发送出去
	strs := strings.Split(data.Content, " ")
	input := strs[len(strs)-1]
	res := "无法识别该操作"
	println("输入------------------------" + input + "------------------------")
	if strings.Contains(data.Content, "> /创建笔记") { //输入笔记标题创建笔记
		if strings.Trim(input, " ") == "" {
			res = "笔记标题为空，创建失败！"
		}
		res = utils.CreateNote(input, data.Author.ID, data.Author.Username)
	}
	if strings.Contains(data.Content, "> /编写笔记") { //输入内容开始编写已创建的笔记
		if strings.Trim(input, " ") == "" {
			res = "笔记内容为空，编写失败！"
		}
		res = utils.AddNote(input, data.Author.ID, data.Author.Username)
	}
	if strings.Contains(data.Content, "> /修改笔记") { //输入标题修改对应笔记
		if strings.Trim(input, " ") == "" {
			res = "笔记标题为空，修改失败！"
		}
		res = utils.UpdateNote(input, data.Author.ID, data.Author.Username)
	}
	if strings.Contains(data.Content, "> /删除笔记") { //输入笔记标题删除对应笔记
		if strings.Trim(input, " ") == "" {
			res = "笔记标题为空，删除失败！"
		}
		res = utils.RemoveNote(input, data.Author.ID, data.Author.Username, false)
	}
	if strings.Contains(data.Content, "> /查看所有笔记") { //查看所有已创建的笔记标题
		res = utils.GetAllNote()
	}
	if strings.Contains(data.Content, "> /查看笔记") { //输入标题查看笔记
		if strings.Trim(input, " ") == "" {
			res = "笔记标题为空，查看失败！"
		}
		res = utils.GetNote(input)
	}
	if strings.Contains(data.Content, "> /当前打开笔记") { //查看所有已创建的笔记标题
		res = utils.GetCurNote(data.Author.ID, data.Author.Username)
	}
	println("输出------------------------" + res + "------------------------")
	api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Content: res})
	return nil
}

// 以下为开发中学习与调试用代码，以注释

// atMessageEventHandler 处理 @机器人 的消息
// func atMessageEventHandler(event *dto.WSPayload, data *dto.WSATMessageData) error {
// 	if strings.HasSuffix(data.Content, "> hello") { // 如果@机器人并输入 hello 则回复 你好。
// 		api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Content: "你好"})
// 	}
// 	if strings.HasSuffix(data.Content, "> 深圳") { // 如果@机器人并输入 深圳 则回复 深圳天气。
// 		weatherData := utils.GetWeatherByCity("深圳")
// 		api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID,
// 			Content: weatherData.ResultData.CityNm + " " + weatherData.ResultData.Weather + " " + weatherData.ResultData.Days + " " + weatherData.ResultData.Week,
// 			Image:   weatherData.ResultData.WeatherIcon, //天气图片
// 		})
// 	}
// 	if strings.HasSuffix(data.Content, "> /天气") { // 如果@机器人并输入 /天气 城市名 则回复 城市天气。
// 		strs := strings.Split(data.Content, " ")
// 		weatherData := utils.GetWeatherByCity(strs[len(strs)-1])
// 		// api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Ark: utils.CreateArkForTemplate23(weatherData)})
// 		api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID,
// 			Content: weatherData.ResultData.CityNm + " " + weatherData.ResultData.Weather + " " + weatherData.ResultData.Days + " " + weatherData.ResultData.Week,
// 			Image:   weatherData.ResultData.WeatherIcon, //天气图片
// 		})
// 	}
// 	if strings.HasSuffix(data.Content, "> admin") { // 如果@机器人并输入 姓名 则回复 用户。
// 		userList := utils.AddUser("admin", 18)
// 		userStr, _ := json.Marshal(userList)
// 		api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Content: string(userStr)})
// 	}

// 	if strings.Contains(data.Content, "> /天气") {
// 		strs := strings.Split(data.Content, " ")
// 		//获取深圳的天气数据
// 		weatherData := utils.GetWeatherByCity(strs[len(strs)-1])
// 		api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Ark: utils.CreateArkForTemplate23(weatherData)})
// 	}
// 	if strings.Contains(data.Content, "> 用户") {
// 		user := getUser()
// 		api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Content: user.Username + " " + data.Author.Username})
// 	}
// 	return nil
// }

// 获取用户
// func getUser() *dto.User {
// 	token := token.BotToken(config.AppID, config.Token)
// 	api := botgo.NewOpenAPI(token).WithTimeout(3 * time.Second)
// 	ctx := context.Background()

// 	user, meError := api.Me(ctx)
// 	if meError != nil {
// 		log.Fatalln("调用 Me 接口失败, err = ", meError)
// 	}
// 	return user
// }
