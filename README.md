# raft-example

Simple key-value implementation of raft

## Usage

```bash
go build
```

Start node-1 as a leader:

```bash
./raft-example --node-id node-1 --raft-port 4321 --http-port 3000
```

Start node-2 as a follower:

```bash
./raft-example --node-id node-2 --raft-port 4322 --http-port 3001
```

To add node-2 as a follower you should call `join` endpoint on node-1:

```bash
curl 'localhost:3000/join?followerAddress=localhost:4322&followerId=node-2'
```

To add value to the cluster you should call `set` endpoint on node-1:

```bash
curl -X POST 'localhost:3000/set' -d '{"key": "x", "value": "23"}' -H 'content-type: application/json'
```

To get value from the cluster you should call `get` endpoint on node-1:

```bash
curl 'localhost:3000/get?key=x'
```

