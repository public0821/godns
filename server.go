package main

import (
    //"./dns"
    "bufio"
    "io"
    "log"
    "math/rand"
    "net"
    "os"
    "strings"
    "time"
)

const FORWARD_PORT_NUM = 5
const MAX_BUF_LEN = 65535

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

var gsession Session
var grecordManager RecordManager

func doForwardToResolver(server *net.UDPConn, forwardConns []*net.UDPConn, resolverAddrs []string, recvMsgChannel chan RecvMsg) {
    buf := make([]byte, MAX_BUF_LEN)
    log.Println("start doForwardToResolver")
    addrIndex := 0
    connIndex := 0
    rand.Seed(time.Now().UTC().UnixNano())
    for recvMsg := range recvMsgChannel {
        //unpack message
        var msg Message
        _, err := msg.UnpackHeaderAndQuestion(recvMsg.data)
        if err != nil {
            log.Println(err)
            continue
        }
        //check whether the domain name is in record manager
        question := &msg.Question[0]
        var record Record
        record.Name = question.Name
        record.Class = question.Type
        record.Type = question.Type
        records, err := grecordManager.QueryRecord(&record)
        if err != nil {
            log.Println(err)
            continue
        }
        //construct an answer record and send to client
        if len(records) > 0 {
            rr, err := RRConstruct(&records[0])
            if err != nil {
                log.Println(err)
                continue
            }
            msg.Hdr.Rcode = RCODE_SUCCESS
            msg.Hdr.QueryResponse = QR_RESPONSE
            msg.Answer = append(msg.Answer, rr)
            buflen, err := msg.Pack(buf, true)
            log.Println("construct an answer record")
            _, err = server.WriteToUDP(buf[:buflen], &recvMsg.addr)
            if err != nil {
                log.Println(err)
                continue
            }
            continue
        }
        //forward to upstream resolver
        addr, err := net.ResolveUDPAddr("udp", resolverAddrs[addrIndex])
        if err != nil {
            log.Println(err)
            continue
        }
        log.Println("write request to resolver")
        conn := forwardConns[connIndex]
        //modify id
        newId := uint16(rand.Int())
        recvMsg.data[0] = uint8(newId >> 8)
        recvMsg.data[1] = uint8(newId)
        _, err = conn.WriteToUDP(recvMsg.data, addr)
        if err != nil {
            log.Println(err)
            continue
        }
        connIndex += 1
        if connIndex >= len(forwardConns) {
            connIndex = 0
        }
        addrIndex += 1
        if addrIndex >= len(resolverAddrs) {
            addrIndex = 0
        }

        var key SessionKey
        key.forwardConn = conn
        key.name = question.Name
        key.id = newId

        var value SessionValue
        value.clientAddr = recvMsg.addr
        value.id = msg.Hdr.Id

        gsession.Lock()
        gsession.buffer[key] = value
        gsession.Unlock()
    }
}

func doRecvFromResolver(server, conn *net.UDPConn) {
    log.Println("start  doRecvFromResolver")
    buf := make([]byte, MAX_BUF_LEN)
    for {
        length, _, err := conn.ReadFromUDP(buf)
        if err != nil {
            log.Println(err)
            continue
        }
        log.Println("receive response from resolver")
        //unpack message
        var msg Message
        _, err = msg.UnpackHeaderAndQuestion(buf[:length])
        if err != nil {
            log.Println(err)
            continue
        }

        var key SessionKey
        key.forwardConn = conn
        question := &msg.Question[0]
        key.name = question.Name
        key.id = msg.Hdr.Id

        gsession.Lock()
        value, ok := gsession.buffer[key]
        if ok {
            delete(gsession.buffer, key)
            gsession.Unlock()
        } else {
            log.Println("key not find in session: ", key)
            gsession.Unlock()
            continue
        }

        log.Println("write response to client")
        //modify id
        newId := value.id
        buf[0] = uint8(newId >> 8)
        buf[1] = uint8(newId)
        _, err = server.WriteToUDP(buf[:length], &value.clientAddr)
        if err != nil {
            log.Println(err)
            continue
        }
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
        log.Println(err)
        return
    }
    var resolverAddrs []string
    for _, dnsServer := range dnsServers {
        resolverAddrs = append(resolverAddrs, dnsServer+":53")
    }

    //initialize session
    gsession.buffer = make(map[SessionKey]SessionValue)
    grecordManager = &RecordManagerSqlite3{}
    err = grecordManager.Open()
    if err != nil {
        log.Println(err)
        return
    }
    defer grecordManager.Close()
    var record Record
    record.Name = "www.test.com"
    record.Class = 1
    record.Type = 1
    record.Value = "10.32.171.60"
    err = grecordManager.AddRecord(&record)
    if err != nil {
        log.Println(err)
        return
    }
    //listen for service
    addr := net.UDPAddr{Port: 53}
    server, err := net.ListenUDP("udp", &addr)
    if err != nil {
        log.Println(err)
        return
    }
    defer server.Close()

    //initialize forward conn
    forwardConns := make([]*net.UDPConn, FORWARD_PORT_NUM)
    for i := 0; i < FORWARD_PORT_NUM; i++ {
        forwardConn, err := net.ListenUDP("udp", &net.UDPAddr{})
        if err != nil {
            log.Println(err)
            return
        }
        forwardConns[i] = forwardConn
        go doRecvFromResolver(server, forwardConn)
        defer forwardConn.Close()
    }

    recvMsgChannel := make(chan RecvMsg, 10)
    go doForwardToResolver(server, forwardConns, resolverAddrs, recvMsgChannel)

    //start web
    err = WebStart()
    if err != nil {
        log.Println(err)
        return
    }

    buf := make([]byte, 65535)
    for {
        length, client, err := server.ReadFromUDP(buf)
        if err != nil {
            log.Println(err)
            continue
        }
        //log.Println(buf[:length])
        recvMsgChannel <- RecvMsg{addr: *client, data: buf[:length]}
    }
}
