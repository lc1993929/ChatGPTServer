package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
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

	router := gin.New()
	router.Use(LogMiddleWare(), gin.Recovery())
	// 允许跨域访问
	router.Use(corsMiddleware())
	router.POST("/send", func(context *gin.Context) {

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

	err := router.Run()
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
	if len(msg) > 197 {
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
	body, _ := io.ReadAll(res.Body)

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

// 跨域中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func LogMiddleWare() gin.HandlerFunc {

	return func(c *gin.Context) {

		method := c.Request.Method
		reqUrl := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logrus.Error(err)
			}
		}(c.Request.Body)

		data, _ := io.ReadAll(c.Request.Body)
		body := string(data[:])
		logrus.WithFields(logrus.Fields{
			"method":      method,
			"uri":         reqUrl,
			"status_code": statusCode,
			"client_ip":   clientIP,
			"body":        body,
		}).Info()

		c.Request.Body = io.NopCloser(bytes.NewBuffer(data))

		c.Next()
	}
}
