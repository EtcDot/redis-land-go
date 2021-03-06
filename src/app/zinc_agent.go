package main

import (
    "os"
    "log"
    "net"
    "net/rpc"
    "custom_jsonrpc"
)

type ZincAgent struct {
    listener net.Listener
    addr string
    db  *Leveldb
    quit_chan chan int
}

func StartZincAgent(agent *ZincAgent) {
    rpc.Register(agent)
    tcpAddr, err := net.ResolveTCPAddr("tcp", agent.addr)
    if err != nil {
        log.Printf("Error: %v", err)
        os.Exit(1)
    }
    listener, err := net.ListenTCP("tcp", tcpAddr)
    if err != nil {
        log.Printf("Error: %v", err)
        os.Exit(1)
    }
    log.Printf("Start json rpc on %v", tcpAddr)
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        log.Printf("New conn:%v", conn)
        go custom_jsonrpc.ServeConn(conn)
    }
}

func (agent *ZincAgent) Get(key *string, value *string) error {
    log.Printf("zinc agent get:%v", *key)
    t, _ := agent.db.Get([]byte(*key))
    *value = string(t)
    return nil
}

func NewZincAgent(addr string, db *Leveldb) *ZincAgent {
    return &ZincAgent{nil, addr, db, make(chan int)}
}
