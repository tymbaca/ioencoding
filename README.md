# IO Encoding

Provides basic wrapper for `xml.Encoder` and `json.Encoder` to encode into
`io.ReadCloser` (using `io.Pipe` internally).

Example:

```go
package main

import "github.com/tymbaca/ioencoding/ioxml"

func main() {
    myObj := // ... take your object

    // Put it inside of encoder, similar to stdlib xml/json
    r := ioxml.NewEncoder(myobj).Indent("", "  ").Encode()
    defer r.Close()

    // use r as any io.Reader
}
```
