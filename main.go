package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"sync"
)

func main() {
	db := &sync.Map{}
	sm := &statemachine{db}

	dataDir := "data"
	err := os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Could not create data directory: %s", err)
	}

	rCfg := &raftconfig{}
	r, err := rCfg.boot(path.Join(dataDir, "raft"+cfg.id), "localhost:"+cfg.raftPort, cfg.id, sm)
	if err != nil {
		log.Fatal(err)
	}

	hs := httpServer{r, db}

	http.HandleFunc("/set", hs.setHandler)
	http.HandleFunc("/get", hs.getHandler)
	http.HandleFunc("/join", hs.joinHandler)
	http.ListenAndServe(":"+cfg.httpPort, nil)
}
