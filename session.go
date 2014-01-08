package main

import (
    //"./dns"
    //"bufio"
    //"io"
    //"log"
    "net"
    //"os"
    //"strings"
    "sync"
)

type SessionKey struct {
    name string
    id   uint16
    port uint16
}

type SessionValue struct {
    name string
    id   string
    ip   net.IP
    port uint16
}

type Session struct {
    sync.RWMutex
    buffer map[SessionKey]SessionValue
}
