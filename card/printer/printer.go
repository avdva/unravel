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

// Logger writes logs.
type Logger interface {
	Printf(format string, args ...interface{})
}

// StringHasher returns hash of a string.
type StringHasher func(s string) []byte

// PrintHandler processes card events and prints them to 'out'.
type PrintHandler struct {
	out    io.Writer
	hasher StringHasher
	l      Logger

	m    sync.RWMutex
	data map[string]cardData
}

// New returns new PrintHandler.
func New(out io.Writer, hasher StringHasher, l Logger) *PrintHandler {
	return &PrintHandler{
		out:    out,
		hasher: hasher,
		l:      l,
		data:   make(map[string]cardData),
	}
}

// OnResize is a window resize callback.
func (ph *PrintHandler) OnResize(h card.EventHeader, from, to card.Dimension) {
	data := ph.doResize(h, from, to)
	ph.doPrint(fmt.Sprintf("%+v\n", data))
}

func (ph *PrintHandler) doResize(h card.EventHeader, from, to card.Dimension) cardData {
	ph.m.Lock()
	defer ph.m.Unlock()
	data := ph.ensureData(h)
	data.ResizeFrom = from
	data.ResizeTo = to
	ph.writeData(data)
	return copyData(&data)
}

// OnCopyPaste is a copt/paste event callback.
func (ph *PrintHandler) OnCopyPaste(h card.EventHeader, form string, pasted bool) {
	data := ph.doCopyPaste(h, form, pasted)
	ph.doPrint(fmt.Sprintf("%+v\n", data))
}

func (ph *PrintHandler) doCopyPaste(h card.EventHeader, form string, pasted bool) cardData {
	ph.m.Lock()
	defer ph.m.Unlock()
	data := ph.ensureData(h)
	data.CopyAndPaste[form] = pasted
	ph.writeData(data)
	return copyData(&data)
}

// OnSubmit is a submit button callback.
func (ph *PrintHandler) OnSubmit(h card.EventHeader, time int) {
	data := ph.doSubmit(h, time)
	ph.doPrint(fmt.Sprintf("%+v, hash = %X\n", data, ph.hasher(data.WebsiteUrl)))
}

func (ph *PrintHandler) doSubmit(h card.EventHeader, time int) cardData {
	ph.m.Lock()
	defer ph.m.Unlock()
	data := ph.ensureData(h)
	data.FormCompletionTime = time
	ph.writeData(data)
	return copyData(&data)
}

func (ph *PrintHandler) doPrint(s string) {
	if _, err := ph.out.Write([]byte(s)); err != nil {
		if ph.l != nil {
			ph.l.Printf("ph: failed to write: %v", err)
		}
	}
}

// ensureData returns existing, or creates new cardData.
// needs ph.m to be Locked.
func (ph *PrintHandler) ensureData(h card.EventHeader) cardData {
	data, found := ph.data[h.SessionID]
	if !found {
		data.SessionId = h.SessionID
		data.WebsiteUrl = h.WebsiteURL
		data.CopyAndPaste = make(map[string]bool)
	}
	return data
}

// writeData stores cardData.
// needs ph.m to be Locked.
func (ph *PrintHandler) writeData(data cardData) {
	ph.data[data.SessionId] = data
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
