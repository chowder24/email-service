package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	mailjet "github.com/mailjet/mailjet-apiv3-go/v3"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

// Sender is a helper for sending and receiving messages.
type Sender struct {
	client *redis.Client
}

type AppConfig struct {
	Publickey     string `json:"MJ_APIKEY_PUBLIC"`
	PrivateKey    string `json:"MJ_APIKEY_PRIVATE"`
	RedisAddr     string  `json:"RedisAddr"`
	RedisPassword  string `json:"RedisPassword"`
}

var(
	config = ReadConf("email-config.json")
)

// New is a constructor for Sender.
func New() Sender {
	return Sender{
		redis.NewClient(&redis.Options{
			Addr:    config.RedisAddr,
			Password: config.RedisPassword,
			DB:       0,
		}),
	}
}


func (b Sender) receiveTask(key string) string {
	result, _ := b.client.RPop(key).Result()
	fmt.Println(result)
	return result
}

func handleTask(value string) {

	delimiter := "|"

	email := strings.Split(value, delimiter)[0]
	name := strings.Join(strings.Split(value, delimiter)[1:], delimiter)

	apiKey := config.Publickey
	secretKey := config.PrivateKey
	mj := mailjet.NewMailjetClient(apiKey, secretKey)
	messagesInfo := []mailjet.InfoMessagesV31 {
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: "akinmolayanoluwatoni@gmail.com",
				Name: "Chowder Admin",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31 {
					Email: email,
					Name: name,
				},
			},
			Subject: "Welcome to Chowder!",
			TextPart: "Dear "+name+" , welcome to Chowder! May the delivery force be with you!",
			HTMLPart: "<h3>Dear "+name+", welcome to <a href=\"https://www.mailjet.com/\">Mailjet</a>!</h3><br />May the delivery force be with you!",
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo }
	res, err := mj.SendMailV31(&messages)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Data: %+v\n", res)
}

func main() {
	TaskKey := "email:new"
	sender := New()
	for {
		time.Sleep(time.Second)
		taskValue := sender.receiveTask(TaskKey)
		if len(taskValue) > 0 {
			go handleTask(taskValue)
		}
	}
}


func ReadConf(conf string) AppConfig {
	data, err := ioutil.ReadFile(conf)
	if err != nil {
		panic(err)
	}
	obj := AppConfig{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		panic(err)
	}
	return obj
}

