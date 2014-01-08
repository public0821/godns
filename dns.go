package main

import (
    //"./dns"
    "bufio"
    "io"
    "log"
    "net"
    "os"
    "strings"
)

const FORWARD_PORT_NUM = 5

func getDnsServer() (servers []string, err error) {
    file, err := os.Open("/etc/resolv.conf")
    if err != nil {
        return
    }
    defer file.Close()
    reader := bufio.NewReader(file)
    lineBytes, _, err := reader.ReadLine()
    for err == nil {
        line := string(lineBytes)
        //remove comments
        commentIndex := strings.Index(line, "#")
        if commentIndex != -1 {
            line = line[:commentIndex]
        }

        fields := strings.Fields(line)
        if len(fields) == 2 && fields[0] == "nameserver" {
            servers = append(servers, fields[1])
        }
        lineBytes, _, err = reader.ReadLine()
    }
    if err == io.EOF {
        err = nil
    }
    return
}

func doForwardToResolver(resolverAddrs []string) {
    _, err = remote.WriteToUDP(buf[:length], &remoteAddr)
    if err != nil {
        log.Panicln(err)
        continue
    }
}
func doReturnToClient() {
    _, err = conn.WriteToUDP(buf[:length], client)
    if err != nil {
        log.Panicln(err)
        continue
    }
}
func doRecvFromResolver(conn *net.UDPConn) {
    length, _, err = remote.ReadFromUDP(buf)
    if err != nil {
        log.Panicln(err)
        continue
    }
}

type RecvMsg struct {
    addr net.UDPAddr
    data []byte
}

func main() {
    //initialize log
    log.SetOutput(os.Stdout)
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    //initialize resolver addrs
    dnsServers, err := getDnsServer()
    if err != nil {
        log.Panicln(err)
        return
    }
    var resolverAddrs []string
    for _, dnsServer := range dnsServers {
        resolverAddrs = append(resolverAddrs, dnsServer+":53")
    }

    //listen for service
    addr := net.UDPAddr{Port: 53}
    server, err := net.ListenUDP("udp", &addr)
    if err != nil {
        log.Panicln(err)
        return
    }
    defer server.Close()

    //initialize forward conn
    forwardConns := make([]*net.UDPConn, FORWARD_PORT_NUM)
    for i := 0; i < FORWARD_PORT_NUM; i++ {
        forwardConn, err := net.ListenUDP("udp", &net.UDPAddr{})
        if err != nil {
            log.Panicln(err)
            return
        }
        forwardConns[i] = forwardConn
        go doRecvFromResolver(forwardConn)
        defer forwardConn.Close()
    }

    go doForwardToResolver(resolverAddrs)

    recvMsgChannel := make(chan RecvMsg, 10)
    buf := make([]byte, 65535)
    for {
        length, client, err := conn.ReadFromUDP(buf)
        if err != nil {
            log.Println(err)
            continue
        }
        recvMsgChannel <- RecvMsg{addr: client, data: buf[:length]}
    }
}
