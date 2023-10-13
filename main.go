package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ijt/go-anytime"
)

type slackStatus struct {
	Text   string `json:"status_text,omitempty"`
	Emoji  string `json:"status_emoji,omitempty"`
	Expiry int    `json:"status_expiration,omitempty"`
}

func (s slackStatus) set(userToken string) {

	// HAX
	type profile struct {
		Profile slackStatus `json:"profile"`
	}

	var p profile

	/* if s.Text == "" && s.Emoji == "" {
		p = profile{Profile: slackStatus{Text: "", Emoji: ""}}
	} */
	p = profile{Profile: s}
	// end HAX

	//fmt.Printf("%#v", p)

	jsonBody, _ := json.Marshal(p)
	bodyReader := bytes.NewReader(jsonBody)

	url := "https://slack.com/api/users.profile.set"
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", "Bearer "+userToken)

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(res)

	//_, err = io.ReadAll(res.Body)
	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error reading HTTP response body: %v", err)
	}

	if res.StatusCode == 200 {
		fmt.Println("ok")
	} else {
		log.Println("We got the response:", string(responseBytes))
	}

}

func main() {

	userToken := os.Getenv("USER_TOKEN")
	if userToken == "" {
		fmt.Println("Put your user token in your env as USER_TOKEN. https://api.slack.com/authentication/token-types#user")
		return
	}

	if len(os.Args) != 4 {
		fmt.Println("takes args, status_text, status_emoji, and optional status_expiry space seperated .. eg:")
		fmt.Println("./main 'this is a text status' ':ant:' 'monday 9am' # set text, emoji until monday 9am")
		fmt.Println("./main 'this is a text status' '' '' # set just a text status, immediately")
		fmt.Println("./main '' '' '' # remove your current status, CURRENTLY BROKEN")
		return
	}

	text := os.Args[1]
	emoji := os.Args[2]

	status := slackStatus{
		Text:  text,
		Emoji: emoji,
	}

	// expiry is tricky ..
	if os.Args[3] != "" {
		expiryDate, err := anytime.Parse(os.Args[3], time.Now())

		if err != nil {
			fmt.Println(err)
			return
		}

		status.Expiry = int(expiryDate.Unix())

	}
	//fmt.Println(expiryDate)

	status.set(userToken)

}
