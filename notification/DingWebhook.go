package notification

import (
	"horizon/config"
	"io"
	"net/http"
	"strings"
	"vitess.io/vitess/go/vt/log"
)

type DingContentMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
type DingContentAt struct {
	AtMobiles []string `json:"atMobiles"`
	AtUserIds []string `json:"atUserIds"`
	IsAtAll   bool     `json:"isAtAll"`
}
type DingContent struct {
	MsgType         string              `json:"msgtype"`
	ContentMarkdown DingContentMarkdown `json:"markdown"`
	ContentAt       DingContentAt       `json:"at"`
}

func SendDingDing(content string) {
	//创建一个请求
	req, err := http.NewRequest("POST", config.Conf.DingWebhook.Webhook, strings.NewReader(content))
	if err != nil {
		log.Errorf("create request failed, err: %v", err.Error())
	}
	client := &http.Client{}
	//设置请求头
	req.Header.Set("Content-Type", "application/json")
	//发送请求
	resp, err := client.Do(req)
	//关闭请求
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("response close failed, err: %v", err.Error())
		}
	}(resp.Body)

	if err != nil {
		log.Errorf("request send ding msg failed, err: %v", err.Error())
	}
}
