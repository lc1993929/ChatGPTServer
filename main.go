package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Query struct {
	Msg string `json:"msg"`
}

func main() {
	r := gin.Default()
	r.POST("/send", func(context *gin.Context) {

		query := Query{}

		err := context.ShouldBindJSON(&query)
		if err != nil {
			log.Panic(err)
		}
		log.Println(query)

		back := sendChatGPT(query.Msg)
		log.Println(back)
		context.JSON(http.StatusOK, gin.H{"msg": back})
	})

	err := r.Run()
	if err != nil {
		log.Panic(err)
	}
}

func sendChatGPT(msg string) string {
	if len(msg) > 97 {
		log.Panic("消息内容长度不能大于197个字节")
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
	req.Header.Add("Authorization", "Bearer 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") //替换成你的API KEY

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Panic(err)
		}
	}(res.Body)
	body, _ := ioutil.ReadAll(res.Body)

	var data map[string]interface{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Panic(err)
	}

	output := data["choices"].([]interface{})[0].(map[string]interface{})["text"].(string)
	return output
}
