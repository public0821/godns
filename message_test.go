package main

import (
    //"fmt"
    "net"
    "testing"
)

func TestPackUnpack(t *testing.T) {
    testPackUnpack(t, false)
    testPackUnpack(t, true)
}

func testPackUnpack(t *testing.T, compression bool) {
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
    //fmt.Println(msg.String())
    buf := make([]byte, 1024)
    length, err := msg.Pack(buf, compression)
    if err != nil {
        t.Error(err)
        return
    }
    //fmt.Println(msg)
    //fmt.Println(length, buf[:length])
    var newMsg Message
    err = newMsg.UnpackAll(buf[:length])
    if err != nil {
        t.Error(err)
        return
    }
    //fmt.Println(newMsg)
    if msg.Hdr != newMsg.Hdr {
        t.Error("msg.Hdr != newMsg.Hdr")
    }
    for i, q := range msg.Question {
        if q != newMsg.Question[i] {
            //fmt.Println(q)
            //fmt.Println(newMsg.Question[i])
            t.Error("msg.Question != newMsg.Question")
        }
    }
    for i, rr := range msg.Answer {
        if rr.String() != newMsg.Answer[i].String() {
            t.Error("msg.Answer != newMsg.Answer")
        }
    }
    for i, rr := range msg.Authority {
        if rr.String() != newMsg.Authority[i].String() {
            t.Error("msg.Authority != newMsg.Authority")
        }
    }
    for i, rr := range msg.Additional {
        if rr.String() != newMsg.Additional[i].String() {
            t.Error("msg.Additional != newMsg.Additional")
        }
    }
}
