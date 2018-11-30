// Copyright 2018 Aleksandr Demakin. All rights reserved.

package printer

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/avdva/unravel/card"
	"github.com/avdva/unravel/hash"

	"github.com/stretchr/testify/require"
)

type errorWriter struct{}

func (e *errorWriter) Write(buff []byte) (int, error) {
	return 0, errors.New("errorWriter")
}

type logger struct {
	e error
}

func (l *logger) Printf(s string, args ...interface{}) {
	l.e = args[0].(error)
}

func TestPrinter(t *testing.T) {
	r := require.New(t)
	buf := bytes.NewBuffer(nil)
	p := New(buf, hash.MakeHasher("pjw"), nil)

	h := card.EventHeader{
		SessionID:  "1",
		WebsiteURL: "google.com",
	}
	p.OnResize(h, card.Dimension{800, 600}, card.Dimension{1920, 1018})
	bytes, err := ioutil.ReadAll(buf)
	r.NoError(err)
	r.Equal([]byte("{WebsiteUrl:google.com SessionId:1 ResizeFrom:{Width:800 Height:600} ResizeTo:{Width:1920 Height:1018} CopyAndPaste:map[] FormCompletionTime:0}\n"), bytes)

	p.OnCopyPaste(h, "form1", false)
	bytes, err = ioutil.ReadAll(buf)
	r.NoError(err)
	r.Equal([]byte("{WebsiteUrl:google.com SessionId:1 ResizeFrom:{Width:800 Height:600} ResizeTo:{Width:1920 Height:1018} CopyAndPaste:map[form1:false] FormCompletionTime:0}\n"), bytes)

	p.OnCopyPaste(h, "form1", true)
	bytes, err = ioutil.ReadAll(buf)
	r.NoError(err)
	r.Equal([]byte("{WebsiteUrl:google.com SessionId:1 ResizeFrom:{Width:800 Height:600} ResizeTo:{Width:1920 Height:1018} CopyAndPaste:map[form1:true] FormCompletionTime:0}\n"), bytes)

	p.OnSubmit(h, 1313)
	bytes, err = ioutil.ReadAll(buf)
	r.NoError(err)
	r.Equal([]byte("{WebsiteUrl:google.com SessionId:1 ResizeFrom:{Width:800 Height:600} ResizeTo:{Width:1920 Height:1018} CopyAndPaste:map[form1:true] FormCompletionTime:1313}, hash = 0E22B00D\n"), bytes)

	h.SessionID = "2"
	h.WebsiteURL = "stackoverflow.com"

	p.OnResize(h, card.Dimension{1920, 1018}, card.Dimension{800, 600})
	bytes, err = ioutil.ReadAll(buf)
	r.NoError(err)
	r.Equal([]byte("{WebsiteUrl:stackoverflow.com SessionId:2 ResizeFrom:{Width:1920 Height:1018} ResizeTo:{Width:800 Height:600} CopyAndPaste:map[] FormCompletionTime:0}\n"), bytes)

	p.OnSubmit(h, 123)
	bytes, err = ioutil.ReadAll(buf)
	r.NoError(err)
	r.Equal([]byte("{WebsiteUrl:stackoverflow.com SessionId:2 ResizeFrom:{Width:1920 Height:1018} ResizeTo:{Width:800 Height:600} CopyAndPaste:map[] FormCompletionTime:123}, hash = 0180272D\n"), bytes)

	l := &logger{}
	p = New(&errorWriter{}, hash.MakeHasher("pjw"), l)
	p.OnSubmit(h, 123)
	r.Equal("errorWriter", l.e.Error())
}
