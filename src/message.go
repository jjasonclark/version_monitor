package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type attachment struct {
	Fallback  string `json:"fallback,omitempty"`
	Color     string `json:"color,omitempty"`
	Pretext   string `json:"pretext,omitempty"`
	Title     string `json:"title,omitempty"`
	TitleLink string `json:"title_link,omitempty"`
	Text      string `json:"text,omitempty"`
}

type slackMessage struct {
	Text        string       `json:"text"`
	Attachments []attachment `json:"attachments,omitempty"`
}

func (cr checkResult) output() (string, error) {
	commitLink := fmt.Sprintf(cr.Verify, cr.Result)
	m := slackMessage{
		Attachments: []attachment{
			attachment{
				Color: "good",
				Text:  fmt.Sprintf("%s deployed: %s", cr.Name, commitLink),
			},
		},
	}
	var buf []byte
	b := bytes.NewBuffer(buf)
	if err := json.NewEncoder(b).Encode(&m); err != nil {
		return "", err
	}
	return b.String(), nil
}
