package dns

import (
	"errors"
	"net"
)

const (
	QR_REQUEST  uint16 = 0
	QR_RESPONSE uint16 = 1
)

//                              1  1  1  1  1  1
//0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|                      ID                       |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|                    QDCOUNT                    |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|                    ANCOUNT                    |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|                    NSCOUNT                    |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|                    ARCOUNT                    |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
const HEADER_LENGTH = 12

type Header struct {
	Id                 uint16
	QueryResponse      uint8 //request or response
	Opcode             uint8
	AuthAnswer         bool
	Truncated          bool
	RecursionDesired   bool
	RecursionAvailable bool
    Zero uint8
	Rcode              uint8
	QDCount            uint16
	ANCount            uint16
	NSCount            uint16
	ARCount            uint16
}

//                              1  1  1  1  1  1
//0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|                                               |
///                     QNAME                     /
///                                               /
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|                     QTYPE                     |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|                     QCLASS                    |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
type Question struct {
	Name   string `dns:"cdomain-name"` // "cdomain-name" specifies encoding (and may be compressed)
	QType  uint16
	QClass uint16
}

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

type RRBase struct {
	Name     string `dns:"cdomain-name"`
	Type     uint16
	Class    uint16
	Ttl      uint32
	RdLength uint16 // length of data after header
}

type RR interface {
}

type A struct {
	RRBase
	IPv4 net.IP
}

//+---------------------+
//|        Header       |
//+---------------------+
//|       Question      | the question for the name server
//+---------------------+
//|        Answer       | RRs answering the question
//+---------------------+
//|      Authority      | RRs pointing toward an authority
//+---------------------+
//|      Additional     | RRs holding additional information
//+---------------------+
type Message struct {
	Hdr        Header
	Question   []Question // Holds the RR(s) of the question section.
	Answer     []RR       // Holds the RR(s) of the answer section.
	Authority  []RR       // Holds the RR(s) of the authority section.
	Additional []RR       // Holds the RR(s) of the additional section.
}

// Unpack a binary message to a Msg structure.
func (msg *Message) Unpack(data []byte) (err error) {
	if len(data) < HEADER_LENGTH {
		err = errors.New("message data too short")
        return
	}
    const (
        _QUREY_RESPONSE = 0x80
        _OPCODE = 0x78
        _AUTH_ANSWER = 0x04 
        _TRUNCATED = 0x02
        _RECURSION_DESIRED = 0x01
        _RECURSION_AVAILABLE = 0x80
        _ZERO = 0x70
        _RCODE = 0x08
    )

    //unpack message hearder
	offset := 0
    msg.Hdr.Id = uint16(data[offset])<<8 | uint16(data[offset+1])
    offset += 2
    msg.Hdr.QueryResponse = uint8(data[offset]) & _QUREY_RESPONSE
    msg.Hdr.Opcode =uint8(data[offset]) & _OPCODE
    msg.Hdr.AuthAnswer = uint8(data[offset]) & _AUTH_ANSWER == 1 
    msg.Hdr.Truncated =uint8(data[offset]) & _TRUNCATED == 1
    msg.Hdr.RecursionDesired =uint8(data[offset]) & _RECURSION_DESIRED == 1
    offset += 1
    msg.Hdr.RecursionAvailable =uint8(data[offset]) & _RECURSION_AVAILABLE == 1
    msg.Hdr.Zero =uint8(data[offset]) & _ZERO
    msg.Hdr.Rcode =uint8(data[offset]) & _RCODE
    offset += 1
    msg.Hdr.QDCount = uint16(data[offset])<<8 | uint16(data[offset+1])
    offset += 2
    msg.Hdr.ANCount = uint16(data[offset])<<8 | uint16(data[offset+1])
    offset += 2
    msg.Hdr.NSCount = uint16(data[offset])<<8 | uint16(data[offset+1])
    offset += 2
    msg.Hdr.ARCount = uint16(data[offset])<<8 | uint16(data[offset+1])
    offset += 2

	// Arrays.
    msg.Question = make([]Question, msg.Hdr.QDCount)
	msg.Answer = make([]RR, msg.Hdr.ANCount)
	msg.Authority = make([]RR, msg.Hdr.NSCount)
	msg.Additional = make([]RR, msg.Hdr.ARCount)

	for i := uint16(0); i < msg.Hdr.QDCount; i++ {
		offset, err = unpackQuestion(&msg.Question[i], offset)
		if err != nil {
			return err
		}
	}
	for i := uint16(0); i < msg.Hdr.ANCount; i++ {
		msg.Answer[i], offset, err = unpackRR(data, offset)
		if err != nil {
			return err
		}
	}
	for i := uint16(0); i < msg.Hdr.NSCount; i++ {
		msg.Authority[i], offset, err = unpackRR(data, offset)
		if err != nil {
			return err
		}
	}
	for i := uint16(0); i < msg.Hdr.ARCount; i++ {
		msg.Additional[i], offset, err = unpackRR(data, offset)
		if err != nil {
			return err
		}
	}
	if offset != len(data) {
		err = errors.New("message data too long")
        return err
	}
	return nil
}

func unpackQuestion(question *Question, start int)(end int, err error){
    return
}
func unpackRR(data []byte, start int) (rr RR, end int, err error) {
        return 
}
