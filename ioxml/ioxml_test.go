package ioxml

import (
	"bytes"
	"encoding/xml"
	"io"
	"runtime"
	"testing"
	"time"
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

	prefix := ""
	indent := "  "
	r := NewEncoder(rb).Indent(prefix, indent).Encode()

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, r)
	r.Close()

	expect, err := xml.MarshalIndent(rb, prefix, indent)
	if err != nil {
		panic(err)
	}

	assert(t, string(expect), buf.String())
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

	if after != before-1 {
		t.Fatalf("encoder goroutine must exit after reader close, (before %d, after %d)", before, after)
	}
}

func assert[T comparable](t *testing.T, expect, actual T) {
	if expect != actual {
		t.Fatalf("'%v'\n\n    must be equal to\n\n'%v'", actual, expect)
	}
}
