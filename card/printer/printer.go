// Copyright 2018 Aleksandr Demakin. All rights reserved.

package printer

import (
	"hash"
	"io"
	"sync"

	"github.com/avdva/unravel/card"
)

type cardData struct {
	WebsiteUrl         string
	SessionId          string
	ResizeFrom         card.Dimension
	ResizeTo           card.Dimension
	CopyAndPaste       map[string]bool
	FormCompletionTime int
}

type PrintHandler struct {
	out    io.Writer
	hasher hash.Hash32

	m    sync.RWMutex
	data map[string]cardData
}

func (ph *PrintHandler) OnResize(h card.EventHeader, from, to card.Dimension) {

}

func (ph *PrintHandler) OnCopyPaste(h card.EventHeader, form string, pasted bool) {

}

func (ph *PrintHandler) OnSubmit(h card.EventHeader, time int) {

}
