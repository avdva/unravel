// Copyright 2018 Aleksandr Demakin. All rights reserved.

package card

// EventHeader is a common part for all events.
type EventHeader struct {
	WebsiteURL string
	SessionID  string
}

// Dimension is a width/height pair.
type Dimension struct {
	Width, Height int
}

// Handler handles card-related events.
type Handler interface {
	OnResize(h EventHeader, from, to Dimension)
	OnCopyPaste(h EventHeader, form string, pasted bool)
	OnSubmit(h EventHeader, time int)
}
