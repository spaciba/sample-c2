package main

import (
	"github.com/spaciba/sample_c2/c2-server/listeners"
)

type ServerState struct {
	listeners map[string]*listeners.Listener
}

func NewServerState() *ServerState {
	return &ServerState{listeners: make(map[string]*listeners.Listener)}
}