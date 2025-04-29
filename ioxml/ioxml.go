package ioxml

import (
	"encoding/xml"
	"io"
)

type Encoder struct {
	v any

	prefix, indent string
	withHeader     bool
}

func NewEncoder(v any) *Encoder {
	return &Encoder{
		v: v,
	}
}

func (enc *Encoder) Indent(prefix, indent string) *Encoder {
	enc.prefix = prefix
	enc.indent = indent
	return enc
}

func (enc *Encoder) WithHeader(on bool) *Encoder {
	enc.withHeader = on
	return enc
}

func (enc *Encoder) Encode() io.ReadCloser {
	r, w := io.Pipe()

	e := xml.NewEncoder(w)
	e.Indent(enc.prefix, enc.indent)

	go func() {
		if enc.withHeader {
			_, err := w.Write([]byte(xml.Header))
			if err != nil {
				w.CloseWithError(err)
				return
			}
		}

		err := e.Encode(enc.v)
		if err != nil {
			w.CloseWithError(err)
			return
		}
		w.Close()
	}()

	return r
}
