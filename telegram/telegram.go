package telegram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const base = "https://api.telegram.org/"

var (
	Token   string
	Channel string
)

var client = http.Client{
	Transport: &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	},
}

type GetMeResponse struct {
	Success bool `json:"ok"`
	Result  struct {
		Id    int  `json:"id"`
		IsBot bool `json:"is_bot"`
	} `json:"result"`
}

func SendMessage(message string) error {
	qs := "chat_id=%40" + Channel + "&text=" + message
	url := fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage?%v", Token, qs)

	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil
}

func GetMe() error {
	url := fmt.Sprintf("https://api.telegram.org/bot%v/getMe", Token)
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response GetMeResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	fmt.Printf("My id: %v\nAm I a bot? %v\n", response.Result.Id, response.Result.IsBot)
	return nil
}
