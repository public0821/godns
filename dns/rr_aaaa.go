package dns

import (
	"fmt"
	"github.com/public0821/dnserver/errors"
	"net"
	//"strings"
)

type RRAaaa struct {
	rrBase
	IPv6 net.IP
}

func (a *RRAaaa) PackRData(buf []byte, index int) (offset int, err error) {
	copy(buf[index:], a.IPv6.To16())
	offset = index + net.IPv6len
	return
}

func (a *RRAaaa) UnpackRData(buf []byte, index int) (offset int, err error) {
	offset = index
	if len(buf)-offset < net.IPv6len {
		err = errors.New("data too short")
		return
	}
	copy(a.IPv6, buf[offset:offset+16])
	offset += net.IPv6len
	return
}

func (a *RRAaaa) SetRData(data string) (err error) {
	a.IPv6 = net.ParseIP(data)
	if a.IPv6 == nil {
		err = errors.New("malformed ipv6 address: " + data)
	}
	return
}

func (a *RRAaaa) String() (text string) {
	text += a.Hdr.String()
	text += fmt.Sprintf("IPv6:\t%s\n", a.IPv6.String())
	return
}
