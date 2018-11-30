// Copyright 2018 Aleksandr Demakin. All rights reserved.

package card

type EventHeader struct {
	WebsiteUrl string
	SessionID  string
}

type Dimension struct {
	Width, Height int
}

type Handler interface {
	OnResize(h EventHeader, from, to Dimension)
	OnCopyPaste(h EventHeader, form string, pasted bool)
	OnSubmit(h EventHeader, time int)
}
