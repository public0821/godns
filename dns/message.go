package dns

import (
	"errors"
	"net"
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
const (
	IPV4_LEN = 4
	IPV6_LEN = 16
)

type Header struct {
	Id                 uint16
	QueryResponse      uint8 //request or response
	Opcode             uint8
	AuthAnswer         bool
	Truncated          bool
	RecursionDesired   bool
	RecursionAvailable bool
	Zero               uint8
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
	Name  string `dns:"cdomain-name"` // "cdomain-name" specifies encoding (and may be compressed)
	Type  uint16
	Class uint16
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

type RRHeader struct {
	Name     string `dns:"cdomain-name"`
	Type     uint16
	Class    uint16
	Ttl      uint32
	RdLength uint16 // length of data after header
}

type RR interface {
}

type A struct {
	Hdr  RRHeader
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
func (msg *Message) UnpackHeaderAndQuestion(data []byte) (offset int, err error) {
	if len(data) < HEADER_LENGTH {
		err = errors.New("message data too short")
		return
	}
	const (
		_QUREY_RESPONSE      = 0x80
		_OPCODE              = 0x78
		_AUTH_ANSWER         = 0x04
		_TRUNCATED           = 0x02
		_RECURSION_DESIRED   = 0x01
		_RECURSION_AVAILABLE = 0x80
		_ZERO                = 0x70
		_RCODE               = 0x08
	)

	//unpack message hearder
	offset = 0
	msg.Hdr.Id, offset = unpackUint16(data, offset)
	msg.Hdr.QueryResponse = uint8(data[offset]) & _QUREY_RESPONSE
	msg.Hdr.Opcode = uint8(data[offset]) & _OPCODE
	msg.Hdr.AuthAnswer = uint8(data[offset])&_AUTH_ANSWER == 1
	msg.Hdr.Truncated = uint8(data[offset])&_TRUNCATED == 1
	msg.Hdr.RecursionDesired = uint8(data[offset])&_RECURSION_DESIRED == 1
	offset += 1
	msg.Hdr.RecursionAvailable = uint8(data[offset])&_RECURSION_AVAILABLE == 1
	msg.Hdr.Zero = uint8(data[offset]) & _ZERO
	msg.Hdr.Rcode = uint8(data[offset]) & _RCODE
	offset += 1
	msg.Hdr.QDCount, offset = unpackUint16(data, offset)
	msg.Hdr.ANCount, offset = unpackUint16(data, offset)
	msg.Hdr.NSCount, offset = unpackUint16(data, offset)
	msg.Hdr.ARCount, offset = unpackUint16(data, offset)

	// Arrays.
	msg.Question = make([]Question, msg.Hdr.QDCount)

	for i := uint16(0); i < msg.Hdr.QDCount; i++ {
		offset, err = unpackQuestion(data, offset, &msg.Question[i])
		if err != nil {
			return
		}
	}
	return
}

func (msg *Message) UnpackAll(data []byte) (err error) {
	var offset int
	offset, err = msg.UnpackHeaderAndQuestion(data)
	if err != nil {
		return
	}
	msg.Answer = make([]RR, msg.Hdr.ANCount)
	msg.Authority = make([]RR, msg.Hdr.NSCount)
	msg.Additional = make([]RR, msg.Hdr.ARCount)

	for i := uint16(0); i < msg.Hdr.ANCount; i++ {
		msg.Answer[i], offset, err = unpackRR(data, offset)
		if err != nil {
			return
		}
	}
	for i := uint16(0); i < msg.Hdr.NSCount; i++ {
		msg.Authority[i], offset, err = unpackRR(data, offset)
		if err != nil {
			return
		}
	}
	for i := uint16(0); i < msg.Hdr.ARCount; i++ {
		msg.Additional[i], offset, err = unpackRR(data, offset)
		if err != nil {
			return
		}
	}
	if offset != len(data) {
		err = errors.New("message data too long")
		return
	}
	return nil
}

func unpackUint32(data []byte, index int) (value uint32, offset int) {
	value = uint32(data[index])<<24 | uint32(data[index+1])<<16 | uint32(data[index+2])<<8 | uint32(data[index+3])
	offset = index + 4
	return
}
func unpackUint16(data []byte, index int) (value uint16, offset int) {
	value = uint16(data[index])<<8 | uint16(data[index+1])
	offset = index + 2
	return
}

func unpackDomainName(data []byte, index int) (name string, offset int, err error) {
	dataLen := len(data)
	offset = index
	for {
		if offset+1 > dataLen {
			err = errors.New("out of range")
			return
		}
		labelLen := int(data[offset])
		offset++
		switch labelLen & 0xC0 {
		case 0x00:
			// end of name
			if labelLen == 0x00 {
				if len(name) == 0 {
					name = "."
					return
				} else {
					return
				}
			}
			if offset+labelLen > dataLen {
				err = errors.New("out of range")
				return
			}
			name += string(data[offset : offset+labelLen])
			name += "."
			offset += labelLen
		case 0xC0:
			// pointer to somewhere else in msg.
			// FIXME maybe there's an infinite loop.
			if labelLen > dataLen {
				err = errors.New("ptr out of range")
				return
			}
			tempName, _, tempErr := unpackDomainName(data, labelLen)
			if tempErr != nil {
				return
			}
			name += tempName
			return
		default:
			// 0x80 and 0x40 are reserved
			err = errors.New("fomart error")
			return
		}
	}
	return
}

func unpackQuestion(data []byte, index int, question *Question) (offset int, err error) {
	offset = index
	question.Name, offset, err = unpackDomainName(data, offset)
	if err != nil {
		return
	}
	if offset+4 > len(data) {
		err = errors.New("out of range")
		return
	}
	question.Type, offset = unpackUint16(data, offset)
	question.Class, offset = unpackUint16(data, offset)
	return
}

type RRBase struct {
	Name     string `dns:"cdomain-name"`
	Type     uint16
	Class    uint16
	Ttl      uint32
	RdLength uint16 // length of data after header
}

func unpackRR(data []byte, index int) (rr RR, offset int, err error) {
	offset = index
	var hdr RRHeader
	hdr.Name, offset, err = unpackDomainName(data, offset)
	if err != nil {
		return
	}
	if offset+10 > len(data) {
		err = errors.New("out of range")
		return
	}
	hdr.Type, offset = unpackUint16(data, offset)
	hdr.Class, offset = unpackUint16(data, offset)
	hdr.Ttl, offset = unpackUint32(data, offset)
	hdr.RdLength, offset = unpackUint16(data, offset)
	if hdr.Class != CLASS_INET && hdr.Type != TYPE_A {
		err = errors.New("unimplement")
		return
	}
	if hdr.RdLength != IPV4_LEN {
		err = errors.New("formart error")
		return
	}
	var a A
	a.Hdr = hdr
	a.IPv4 = net.IPv4(data[offset], data[offset+1], data[offset+2], data[offset+3])
	offset += IPV4_LEN
	return
}
