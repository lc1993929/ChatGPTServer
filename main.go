package main

import (
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Query struct {
	Msg string `json:"msg"`
}

func main() {
	configLog()

	apiKey := flag.String("apiKey", "", "openApi apiKey")
	// 解析命令行参数
	flag.Parse()

	logrus.Info(*apiKey)

	r := gin.Default()
	r.POST("/send", func(context *gin.Context) {

		query := Query{}

		err := context.ShouldBindJSON(&query)
		if err != nil {
			logrus.Error(err)
		}

		logrus.Info("send: " + query.Msg)
		back := sendChatGPT(query.Msg, *apiKey)
		logrus.Info("back: " + back)
		context.JSON(http.StatusOK, gin.H{"msg": back})
	})

	err := r.Run()
	if err != nil {
		logrus.Error(err)
	}
}

func configLog() {
	file, err := os.Create("send.log")
	logrus.SetOutput(file)
	if err != nil {
		logrus.Error("Cannot create log file", err)
	}
	gin.DefaultWriter = io.MultiWriter(file)
}

func sendChatGPT(msg string, apiKey string) string {
	if len(msg) > 97 {
		result := "消息内容长度不能大于197个字节"
		logrus.Info(result)
		return result
	}

	url := "https://api.openai.com/v1/completions"

	payload := strings.NewReader(`{
    "model": "text-davinci-003",
    "prompt": "` + msg + `",
	"temperature": 0,
		"max_tokens": 3900
}`)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey) //替换成你的API KEY

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(res.Body)
	body, _ := ioutil.ReadAll(res.Body)

	var data map[string]interface{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		logrus.Error(err)
	}

	if data["error"] != nil {
		return data["error"].(map[string]interface{})["message"].(string)
	}

	output := data["choices"].([]interface{})[0].(map[string]interface{})["text"].(string)
	return output
}
