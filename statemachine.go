package main

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/hashicorp/raft"
)

type statemachine struct {
	db *sync.Map
}

func (sm *statemachine) Apply(log *raft.Log) any {
	switch log.Type {
	case raft.LogCommand:
		var sp command

		if err := sp.Decode(log.Data); err != nil {
			return fmt.Errorf("failed to decode data: %s", err)
		}

		sm.db.Store(sp.Key, sp.Value)

	default:
		return fmt.Errorf("unrecognized command type: %d", log.Type)
	}

	return nil
}

func (sm *statemachine) Restore(rc io.ReadCloser) error {
	sm.db.Range(func(key, _ any) bool {
		sm.db.Delete(key)
		return true
	})

	decoder := json.NewDecoder(rc)

	for decoder.More() {
		var sp command
		if err := decoder.Decode(&sp); err != nil {
			return fmt.Errorf("failed to decode data: %s", err)
		}

		sm.db.Store(sp.Key, sp.Value)
	}

	return rc.Close()
}

func (sm *statemachine) Snapshot() (raft.FSMSnapshot, error) {
	return snapshotNoop{}, nil
}
