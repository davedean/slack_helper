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

	p := profile{Profile: s}
	// end HAX

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

	_, err = io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error reading HTTP response body: %v", err)
	}

	//	log.Println("We got the response:", string(responseBytes))

}

func main() {

	userToken := os.Getenv("USER_TOKEN")

	text := os.Args[1]
	emoji := os.Args[2]

	// parse time
	expiryDate, err := time.Parse("2006-01-02", os.Args[3])

	if err != nil {
		fmt.Println(err)
		return
	}

	status := slackStatus{
		Text:   text,
		Emoji:  emoji,
		Expiry: int(expiryDate.Unix()),
	}

	status.set(userToken)

}
