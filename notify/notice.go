package notify

import (
	"encoding/json"
	"github.com/mounsurf/lib/zhttp"
	"strings"
)

type Notice struct {
	Api     string `yaml:"api"`
	Keyword string `yaml:"keyword"`
}

type Message struct {
	MsgType string      `json:"msgtype"`
	Text    MessageText `json:"text"`
}

type MessageText struct {
	Content string `json:"content"`
}

func (n *Notice) SendMessage(content string) error {
	if n == nil || n.Api == "" || n.Keyword == "" || content == "" {
		return nil
	}
	if !strings.Contains(content, n.Keyword) {
		content = n.Keyword + "\n" + content
	}
	message := Message{
		MsgType: "text",
		Text: MessageText{
			Content: content,
		},
	}
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = zhttp.Post(n.Api, &zhttp.RequestOptions{
		UserAgent: "null",
		JSON:      string(data),
	})
	return err
}
