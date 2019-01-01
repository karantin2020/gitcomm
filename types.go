package gitcomm

import (
	"bytes"
	"fmt"
)

// Message type holds all commit message fields
type Message struct {
	// Type message field
	Type string
	// Subject message field
	Subject string
	// Body message field
	Body string
	// Foot message field
	Foot string
}

func (m Message) String() string {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, "%s: %s\n\n%s\n%s\n", m.Type, m.Subject, m.Body, m.Foot)
	return buf.String()
}
