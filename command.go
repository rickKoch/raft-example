package main

import "encoding/json"

type command struct {
	Key   string
	Value string
}

func (cmd *command) Decode(data []byte) error {
	return json.Unmarshal(data, cmd)
}
