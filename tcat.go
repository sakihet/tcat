package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type Client struct {
	Id     string `json:"client_id"`
	Secret string `json:"client_secret"`
}

var (
	message string
)

func main() {
	app := cli.NewApp()
	app.Version = "0.0.2"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "topicId, t",
			Usage: "typetalk topic id to post to",
		},
		cli.BoolFlag{
			Name:  "configure",
			Usage: "configure tcat",
		},
		cli.BoolFlag{
			Name:  "plain, p",
			Usage: "post message as plain text instead of code blocks",
		},
		cli.StringFlag{
			Name:  "syntax, s",
			Usage: "post code with syntax highlighting",
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.Bool("configure") {
			fmt.Printf("input typetalk api client id: ")
			client_id, _, err := bufio.NewReader(os.Stdin).ReadLine()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("input typetalk api client secret: ")
			client_secret, _, err := bufio.NewReader(os.Stdin).ReadLine()
			if err != nil {
				log.Fatal(err)
			}
			client := Client{Id: string(client_id), Secret: string(client_secret)}
			jsonBytes, err := json.Marshal(client)
			if err != nil {
				fmt.Println("JSON Marshal error:", err)
			}
			fmt.Println(string(jsonBytes))
			homedir := os.Getenv("HOME")
			if homedir == "" {
				log.Fatal(err)
			}
			file, err := os.Create(homedir + `/.tcat`)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			file.Write(([]byte)(jsonBytes))
			os.Exit(0)
		}
		homedir := os.Getenv("HOME")
		bytes, err := ioutil.ReadFile(homedir + "/.tcat")
		if err != nil {
			log.Fatal(err)
		}
		var client Client
		if err := json.Unmarshal(bytes, &client); err != nil {
			log.Fatal(err)
		}
		clientId := client.Id
		clientSecret := client.Secret

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

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			message += scanner.Text() + "\n"
		}

		if !c.Bool("plain") {
			syntax := c.String("syntax")
			message = "```" + syntax + "\n" + message + "```"
		}
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
