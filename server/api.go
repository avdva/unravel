// Copyright 2018 Aleksandr Demakin. All rights reserved.

package server

import (
	"encoding/json"
	"io"
)

const (
	evCopyPaste    = "copyAndPaste"
	evWindowResize = "windowResize"
	evSubmit       = "timeTaken"
)

type request struct {
	EventType string `json:"eventType"`
	URL       string `json:"websiteUrl"`
	SessionID string `json:"sessionId"`

	Pasted     *bool      `json:"pasted"`
	FormID     *string    `json:"formId"`
	ResizeFrom *dimension `json:"resizeFrom"`
	ResizeTo   *dimension `json:"resizeTo"`
	Time       *int       `json:"time"`
}

type dimension struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func decodeRequest(r io.Reader) (*request, error) {
	var req request
	if err := json.NewDecoder(r).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
