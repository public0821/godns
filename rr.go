package main

import (
    "fmt"
    "github.com/public0821/dnserver/db"
    "github.com/public0821/dnserver/errors"
    "net"
    //"strings"
)

//                              1  1  1  1  1  1
//0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|                                               |
///                                               /
///                      NAME                     /
//|                                               |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|                      TYPE                     |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|                     CLASS                     |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|                      TTL                      |
//|                                               |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|                   RDLENGTH                    |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--|
///                     RDATA                     /
///                                               /
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

type RRHeader struct {
    Name     string
    Type     uint16
    Class    uint16
    Ttl      uint32
    RDLength uint16 //Don't set this field manually
}

func (h *RRHeader) String() (text string) {
    text += fmt.Sprintf("Name:\t%s\n", h.Name)
    text += fmt.Sprintf("Type:\t%s\n", TypeToString[h.Type])
    text += fmt.Sprintf("Class:\t%s\n", ClassToString[h.Class])
    text += fmt.Sprintf("Ttl:\t%d\n", h.Ttl)
    text += fmt.Sprintf("RDLength:\t%d\n", h.RDLength)
    return
}

type RR interface {
    Header() *RRHeader
    SetHeader(header *RRHeader)
    PackRData(buf []byte, index int) (offset int, err error)
    UnpackRData(buf []byte, index int) (offset int, err error)
    String() string
    SetRData(data string) (err error)
    //RDataLength() int
}

type rrBase struct {
    Hdr RRHeader
}

func (h *rrBase) Header() (header *RRHeader) {
    return &h.Hdr
}
func (h *rrBase) SetHeader(header *RRHeader) {
    h.Hdr = *header
}

type RRA struct {
    rrBase
    IPv4 net.IP
}

func (a *RRA) PackRData(buf []byte, index int) (offset int, err error) {
    copy(buf[index:], a.IPv4.To4())
    offset = index + IPV4_LEN
    return
}

func (a *RRA) UnpackRData(buf []byte, index int) (offset int, err error) {
    offset = index
    if len(buf)-offset < IPV4_LEN {
        err = errors.New("data too short")
        return
    }
    a.IPv4 = net.IPv4(buf[offset], buf[offset+1], buf[offset+2], buf[offset+3])
    offset += IPV4_LEN
    return
}

func (a *RRA) SetRData(data string) (err error) {
    a.IPv4 = net.ParseIP(data)
    if a.IPv4 == nil {
        err = errors.New("malformed ipv4 address: " + data)
    }
    return
}

func (a *RRA) String() (text string) {
    text += a.Hdr.String()
    text += fmt.Sprintf("IPv4:\t%s\n", a.IPv4.String())
    return
}

func RRNew(rrtype uint16) (rr RR, err error) {
    switch rrtype {
    case TYPE_A:
        rr = new(RRA)
        return
    default:
        err = errors.New("unimplement")
        return
    }
}

func RRConstruct(rrecord *db.RRecord) (rr RR, err error) {
    rr, err = RRNew(rrecord.Type)
    if err != nil {
        return
    }
    var header RRHeader
    header.Name = rrecord.Name
    header.Class = rrecord.Class
    header.Type = rrecord.Type
    header.Ttl = rrecord.Ttl
    rr.SetHeader(&header)
    err = rr.SetRData(rrecord.Value)
    return
}
