package dns

import (
	"fmt"
	"net"
	"testing"
)

func TestTest(t *testing.T) {
	var msg Message
	msg.Hdr.Id = 100
	//answers := make(RR[6])
	var a A
	a.Name = "test"
	a.IPv4 = net.ParseIP("192.168.1.1")
	msg.Answer = append(msg.Answer, a)
	msg.Answer = append(msg.Answer, a)
	fmt.Println(len(msg.Question))
	fmt.Println(cap(msg.Question))
	fmt.Println(msg)
	t.Log("pass")
	t.Log("aaapass")
}
