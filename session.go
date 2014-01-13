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
        name        string
        id          uint16
        forwardConn *net.UDPConn
        resolver    string
}

type SessionValue struct {
        name       string
        id         uint16
        clientAddr net.UDPAddr
}

type Session struct {
        sync.RWMutex
        buffer map[SessionKey]SessionValue
}
