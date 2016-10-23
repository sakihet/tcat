package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/urfave/cli"
)

type Auth struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

var (
	clientId     = os.Getenv("TYPETALK_API_CLIENT_ID")
	clientSecret = os.Getenv("TYPETALK_API_CLIENT_SECRET")
	message      string
)

func main() {
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "topicId, t",
			Usage: "typetalk topic id to post to",
		},
	}
	app.Action = func(c *cli.Context) error {
		if (clientId == "") || (clientSecret == "") {
			log.Fatal("env is missing")
		}
		resp, err := http.PostForm(
			"https://typetalk.in/oauth2/access_token",
			url.Values{
				"client_id":     {clientId},
				"client_secret": {clientSecret},
				"grant_type":    {"client_credentials"},
				"scope":         {"topic.post"}})
		if err != nil {
			log.Fatal("authentication failed")
		}
		var d Auth
		err = json.NewDecoder(resp.Body).Decode(&d)
		message = "```"
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			message += "\n" + scanner.Text()
		}
		message += "\n```"
		topicId := c.String("topicId")
		resp, err = http.PostForm(
			fmt.Sprintf("https://typetalk.in/api/v1/topics/%s", topicId),
			url.Values{
				"access_token": {d.AccessToken},
				"message":      {message}})
		return nil
	}
	app.Run(os.Args)
}
