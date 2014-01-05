package dns

import (
	"net"
)

const (
	RR_REQUEST  uint16 = 0
	RR_RESPONSE uint16 = 1
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
type Header struct {
	Id                 uint16
	QueryResponse      uint8 //request or response
	Opcode             uint8
	AuthAnswer         bool
	Truncated          bool
	RecursionDesired   bool
	RecursionAvailable bool
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
	Qtype  uint16
	Qclass uint16
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
	Rdlength uint16 // length of data after header
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
