package iojson

import (
	"encoding/json"
	"io"
)

type Encoder struct {
	v any

	prefix, indent string
	escapeHTML     bool
}

func NewEncoder(v any) *Encoder {
	return &Encoder{
		v:          v,
		escapeHTML: true,
	}
}

func (enc *Encoder) Indent(prefix, indent string) *Encoder {
	enc.prefix = prefix
	enc.indent = indent
	return enc
}

func (enc *Encoder) EscapeHTML(on bool) *Encoder {
	enc.escapeHTML = on
	return enc
}

func (enc *Encoder) Encode() io.ReadCloser {
	r, w := io.Pipe()

	e := json.NewEncoder(w)
	e.SetIndent(enc.prefix, enc.indent)
	e.SetEscapeHTML(enc.escapeHTML)

	go func() {
		err := e.Encode(enc.v)
		if err != nil {
			w.CloseWithError(err)
			return
		}
		w.Close()
	}()

	return r
}
