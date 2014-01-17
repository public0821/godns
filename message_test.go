package main

import (
    "fmt"
    "net"
    "testing"
)

func TestPackUnpack(t *testing.T) {
    a := &RRA{}
    a.Hdr.Class = CLASS_INET
    a.Hdr.Type = TYPE_A
    a.Hdr.Ttl = 60
    a.Hdr.Name = "www.baidu.com"
    a.IPv4 = net.IPv4(192, 168, 1, 1)

    var msg Message
    msg.Hdr.Id = 1234
    msg.Hdr.Zero = 3
    msg.Hdr.AuthAnswer = true
    msg.Hdr.Opcode = OPCODE_QUERY
    msg.Hdr.QueryResponse = QR_RESPONSE
    msg.Hdr.Rcode = RCODE_SUCCESS
    msg.Hdr.RecursionAvailable = true
    msg.Hdr.RecursionDesired = true
    msg.Hdr.Truncated = false

    var question Question
    question.Class = CLASS_INET
    question.Type = TYPE_A
    question.Name = "www.baidu.com"

    msg.Question = append(msg.Question, question)
    msg.Answer = append(msg.Answer, a)
    fmt.Println(msg.String())
    buf := make([]byte, 1024)
    length, err := msg.Pack(buf, false)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(buf[:length])
    //data := []byte{103, 85, 1, 32, 0, 1, 0, 0, 0, 0, 0, 1, 3, 119, 119, 119, 5, 98, 97, 105
, 100, 117, 3, 99, 111, 109, 0, 0, 1, 0, 1, 0,
    //0, 41, 16, 0, 0, 0, 0, 0, 0, 0}
    //var msg Message
    //fmt.Println(msg.Unpack(data))
    //fmt.Println(msg)
}
