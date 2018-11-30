// Copyright 2018 Aleksandr Demakin. All rights reserved.

package printer

import (
	"fmt"
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

// StringHasher32 returns hash of a string.
type StringHasher32 func(s string) []byte

// PrintHandler processed card events and prints them to 'out'.
type PrintHandler struct {
	out    io.Writer
	hasher StringHasher32

	m    sync.RWMutex
	data map[string]cardData
}

// New returns new PrintHandler.
func New(out io.Writer, hasher StringHasher32) *PrintHandler {
	return &PrintHandler{
		out:    out,
		hasher: hasher,
		data:   make(map[string]cardData),
	}
}

// OnResize is a window resize callback.
func (ph *PrintHandler) OnResize(h card.EventHeader, from, to card.Dimension) {
	data := ph.doResize(h, from, to)
	ph.doPrint(fmt.Sprintln(data))
}

func (ph *PrintHandler) doResize(h card.EventHeader, from, to card.Dimension) cardData {
	data := ph.ensureData(h)
	data.ResizeFrom = from
	data.ResizeTo = to
	ph.m.Lock()
	defer ph.m.Unlock()
	ph.data[h.SessionID] = data
	return copyData(&data)
}

// OnCopyPaste is a copt/paste event callback.
func (ph *PrintHandler) OnCopyPaste(h card.EventHeader, form string, pasted bool) {
	data := ph.doCopyPaste(h, form, pasted)
	ph.doPrint(fmt.Sprintln(data))
}

func (ph *PrintHandler) doCopyPaste(h card.EventHeader, form string, pasted bool) cardData {
	data := ph.ensureData(h)
	data.CopyAndPaste[form] = pasted
	ph.m.Lock()
	defer ph.m.Unlock()
	ph.data[h.SessionID] = data
	return copyData(&data)
}

// OnSubmit is a submit button handler.
func (ph *PrintHandler) OnSubmit(h card.EventHeader, time int) {
	data := ph.doSubmit(h, time)
	ph.doPrint(fmt.Sprintf("%v, hash = %X\n", data, ph.hasher(data.WebsiteUrl)))
}

func (ph *PrintHandler) doSubmit(h card.EventHeader, time int) cardData {
	data := ph.ensureData(h)
	data.FormCompletionTime = time
	ph.m.Lock()
	defer ph.m.Unlock()
	ph.data[h.SessionID] = data
	return copyData(&data)
}

func (ph *PrintHandler) doPrint(s string) {
	ph.out.Write([]byte(s))
}

// ensureData returns existing, or creates new cardData.
func (ph *PrintHandler) ensureData(h card.EventHeader) cardData {
	ph.m.RLock()
	defer ph.m.RUnlock()
	data, found := ph.data[h.SessionID]
	if !found {
		data.SessionId = h.SessionID
		data.WebsiteUrl = h.WebsiteUrl
		data.CopyAndPaste = make(map[string]bool)
	}
	return data
}

// copyData copies the entire structure.
func copyData(data *cardData) cardData {
	result := *data
	result.CopyAndPaste = make(map[string]bool)
	for formID, pasted := range data.CopyAndPaste {
		result.CopyAndPaste[formID] = pasted
	}
	return result
}
