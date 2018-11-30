// Copyright 2018 Aleksandr Demakin. All rights reserved.

package printer

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/avdva/unravel/card"
	"github.com/avdva/unravel/hash"

	"github.com/stretchr/testify/require"
)

func TestPrinter(t *testing.T) {
	r := require.New(t)
	buf := bytes.NewBuffer(nil)
	p := New(buf, hash.MakeHasher("pjw"), nil)

	h := card.EventHeader{
		SessionID:  "1",
		WebsiteUrl: "google.com",
	}
	p.OnResize(h, card.Dimension{800, 600}, card.Dimension{1920, 1018})
	bytes, err := ioutil.ReadAll(buf)
	r.NoError(err)
	r.Equal([]byte("{google.com 1 {800 600} {1920 1018} map[] 0}\n"), bytes)

	p.OnCopyPaste(h, "form1", false)
	bytes, err = ioutil.ReadAll(buf)
	r.NoError(err)
	r.Equal([]byte("{google.com 1 {800 600} {1920 1018} map[form1:false] 0}\n"), bytes)

	p.OnCopyPaste(h, "form1", true)
	bytes, err = ioutil.ReadAll(buf)
	r.NoError(err)
	r.Equal([]byte("{google.com 1 {800 600} {1920 1018} map[form1:true] 0}\n"), bytes)

	p.OnSubmit(h, 1313)
	bytes, err = ioutil.ReadAll(buf)
	r.NoError(err)
	r.Equal([]byte("{google.com 1 {800 600} {1920 1018} map[form1:true] 1313}, hash = 0E22B00D\n"), bytes)

	h.SessionID = "2"
	h.WebsiteUrl = "stackoverflow.com"

	p.OnResize(h, card.Dimension{1920, 1018}, card.Dimension{800, 600})
	bytes, err = ioutil.ReadAll(buf)
	r.NoError(err)
	r.Equal([]byte("{stackoverflow.com 2 {1920 1018} {800 600} map[] 0}\n"), bytes)

	p.OnSubmit(h, 123)
	bytes, err = ioutil.ReadAll(buf)
	r.NoError(err)
	r.Equal([]byte("{stackoverflow.com 2 {1920 1018} {800 600} map[] 123}, hash = 0180272D\n"), bytes)
}
