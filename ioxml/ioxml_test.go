package ioxml

import (
	"bytes"
	"encoding/xml"
	"io"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type RecordBook struct {
	Records []Record `xml:"records"`
}

type Record struct {
	Key string `xml:"key"`
	Val string `xml:"val"`
}

func TestEncoder(t *testing.T) {
	rb := RecordBook{
		Records: []Record{
			{Key: "k1", Val: "v1"},
			{Key: "k2", Val: "v2"},
			{Key: "k3", Val: "v3"},
			{Key: "k4", Val: "v4"},
		},
	}

	gsBefore := runtime.NumGoroutine()

	prefix := ""
	indent := "  "
	r := NewEncoder(rb).Indent(prefix, indent).Encode()
	defer r.Close()

	gsMiddle := runtime.NumGoroutine()
	require.Equal(t, gsBefore+1, gsMiddle, "encoder goroutine must spawn when reader created, (before %d, after %d)", gsBefore, gsMiddle)

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, r)

	gsAfter := runtime.NumGoroutine()
	require.Equal(t, gsBefore, gsAfter, "encoder goroutine must exit after reader close, (before %d, after %d)", gsBefore, gsAfter)

	expect, err := xml.MarshalIndent(rb, prefix, indent)
	if err != nil {
		panic(err)
	}

	require.Equal(t, string(expect), buf.String())
}

func TestEncoder_ReaderClose(t *testing.T) {
	rb := RecordBook{
		Records: []Record{
			{Key: "k1", Val: "v1"},
			{Key: "k2", Val: "v2"},
			{Key: "k3", Val: "v3"},
			{Key: "k4", Val: "v4"},
		},
	}

	prefix := ""
	indent := "  "
	r := NewEncoder(rb).Indent(prefix, indent).Encode()

	before := runtime.NumGoroutine()

	io.CopyN(io.Discard, r, 10)
	r.Close()

	time.Sleep(1 * time.Millisecond)
	after := runtime.NumGoroutine()

	require.Equal(t, before-1, after, "encoder goroutine must exit after reader close, (before %d, after %d)", before, after)
}
