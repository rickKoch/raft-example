package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

type raftconfig struct{}

func (c *raftconfig) boot(dir, raftAddress, nodeId string, sm *statemachine) (*raft.Raft, error) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}

	store, err := c.setStore(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot create bolt store: %s", err)
	}

	snaps, err := c.setSnapshot(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot create snapshot store: %s", err)
	}

	transport, err := c.setTransport(raftAddress)
	if err != nil {
		return nil, fmt.Errorf("cannot create TCP transport: %s", err)
	}

	raftCfg := raft.DefaultConfig()
	raftCfg.LocalID = raft.ServerID(nodeId)

	r, err := raft.NewRaft(raftCfg, sm, store, store, snaps, transport)
	if err != nil {
		return nil, fmt.Errorf("cannot create raft: %s", err)
	}

	r.BootstrapCluster(raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      raft.ServerID(nodeId),
				Address: transport.LocalAddr(),
			},
		},
	})

	return r, nil
}

func (c *raftconfig) setStore(dir string) (*raftboltdb.BoltStore, error) {
	return raftboltdb.NewBoltStore(filepath.Join(dir, "raft.db"))
}

func (c *raftconfig) setSnapshot(dir string) (*raft.FileSnapshotStore, error) {
	return raft.NewFileSnapshotStore(filepath.Join(dir, "snaphost"), 2, os.Stderr)
}

func (c *raftconfig) setTransport(raftAddress string) (*raft.NetworkTransport, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", raftAddress)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve TCP address: %s", err)
	}

	return raft.NewTCPTransport(raftAddress, tcpAddr, 10, 10*time.Second, os.Stderr)
}
