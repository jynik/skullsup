// SPDX License Identifier: MIT
package network

import "fmt"

// Command
const (
	CmdColor     = "color"
	CmdReanimate = "reanimate"
	CmdIncant    = "incant"
)

type Message struct {
	Command string   `json:"cmd"`
	Args    []string `json:"args"`
	Period  int      `json:"period"`
}

func (m *Message) String() string {
	return fmt.Sprintf("{ %s %s (%d ms) }", m.Command, m.Args, m.Period)
}
