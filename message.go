package main

import (
    "fmt"
    "strings"
)

const MAX_COMPRESSION_DEPTH = 5

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
    QDCount            uint16 //Don't set this field manually
    ANCount            uint16 //Don't set this field manually
    NSCount            uint16 //Don't set this field manually
    ARCount            uint16 //Don't set this field manually
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
    Name  string
    Type  uint16
    Class uint16
}

func (q *Question) String() (text string) {
    text += fmt.Sprintf("Name:\t%s\n", q.Name)
    text += fmt.Sprintf("Type:\t%s\n", TypeToString[q.Type])
    text += fmt.Sprintf("Class:\t%s\n", ClassToString[q.Class])
    return
}

//func (a *RRA) RDataLength() int {
//return IPV4_LEN
//}

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

func (msg *Message) String() (text string) {
    text = "---HEADER---\n"
    text += fmt.Sprintf("Id:\t%d\n", msg.Hdr.Id)
    text += fmt.Sprintf("QueryResponse:\t%d\n", msg.Hdr.QueryResponse)
    text += fmt.Sprintf("Opcode:\t%s\n", OpcodeToString[msg.Hdr.Opcode])
    text += fmt.Sprintf("AuthAnswer:\t%t\n", msg.Hdr.AuthAnswer)
    text += fmt.Sprintf("Truncated:\t%t\n", msg.Hdr.Truncated)
    text += fmt.Sprintf("RecursionDesired:\t%t\n", msg.Hdr.RecursionDesired)
    text += fmt.Sprintf("RecursionAvailable:\t%t\n", msg.Hdr.RecursionAvailable)
    text += fmt.Sprintf("Zero:\t%d\n", msg.Hdr.Zero)
    text += fmt.Sprintf("Rcode:\t%s\n", RcodeToString[msg.Hdr.Rcode])
    text += fmt.Sprintf("QDCount:\t%d\n", msg.Hdr.QDCount)
    text += fmt.Sprintf("ANCount:\t%d\n", msg.Hdr.ANCount)
    text += fmt.Sprintf("NSCount:\t%d\n", msg.Hdr.NSCount)
    text += fmt.Sprintf("ARCount:\t%d\n", msg.Hdr.ARCount)
    for index, question := range msg.Question {
        text += fmt.Sprintf("---QUERY(%d)---\n", index+1)
        text += question.String()
    }
    for index, answer := range msg.Answer {
        text += fmt.Sprintf("---ANSWER(%d)---\n", index+1)
        text += answer.String()
    }
    for index, ns := range msg.Authority {
        text += fmt.Sprintf("---AUTHORITH(%d)---\n", index+1)
        text += ns.String()
    }
    for index, additional := range msg.Additional {
        text += fmt.Sprintf("---ADDITIONAL(%d)---\n", index+1)
        text += additional.String()
    }
    return
}

// Unpack a binary message to a Msg structure.
func (msg *Message) UnpackHeaderAndQuestion(data []byte) (offset int, err error) {
    if len(data) < HEADER_LENGTH {
        err = NewError("message data too short")
        return
    }

    //unpack message hearder
    offset = 0
    msg.Hdr.Id, offset = unpackUint16(data, offset)
    msg.Hdr.QueryResponse = (uint8(data[offset]) & _QUREY_RESPONSE) >> 7
    msg.Hdr.Opcode = (uint8(data[offset]) & _OPCODE) >> 3
    msg.Hdr.AuthAnswer = uint8(data[offset])&_AUTH_ANSWER != 0
    msg.Hdr.Truncated = uint8(data[offset])&_TRUNCATED != 0
    msg.Hdr.RecursionDesired = uint8(data[offset])&_RECURSION_DESIRED != 0
    offset += 1
    msg.Hdr.RecursionAvailable = uint8(data[offset])&_RECURSION_AVAILABLE != 0
    msg.Hdr.Zero = (uint8(data[offset]) & _ZERO) >> 4
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
        err = NewError("message data too long")
        return
    }
    return nil
}

func unpackUint32(data []byte, index int) (value uint32, offset int) {
    value = uint32(data[index])<<24 | uint32(data[index+1])<<16 | uint32(data[index+2])<<8 |
        uint32(data[index+3])
    offset = index + 4
    return
}
func unpackUint16(data []byte, index int) (value uint16, offset int) {
    value = uint16(data[index])<<8 | uint16(data[index+1])
    offset = index + 2
    return
}

func unpackDomainName(data []byte, index int, maxDepth int) (name string, offset int, err error) {
    dataLen := len(data)
    offset = index
    for {
        if offset+1 > dataLen {
            err = NewError("out of range")
            return
        }
        labelLen := int(data[offset])
        offset++
        switch uint8(labelLen) & 0xC0 {
        case 0x00:
            // end of name
            if labelLen == 0x00 {
                if len(name) == 0 {
                    name = "."
                    return
                } else {
                    //remove the last dot
                    name = name[:len(name)-1]
                    return
                }
            }
            if offset+labelLen > dataLen {
                err = NewError("out of range")
                return
            }
            name += string(data[offset : offset+labelLen])
            name += "."
            offset += labelLen
        case 0xC0:
            if offset+1 > dataLen {
                err = NewError("out of range")
                return
            }
            lablePtr := uint16(data[offset-1])<<10>>2 | uint16(data[offset])
            offset++
            // pointer to somewhere else in message.
            if int(lablePtr) > dataLen {
                err = NewError("ptr out of range")
                return
            }
            if maxDepth == 0 {
                err = NewError("too many ptr")
                return
            }
            tempName, _, tempErr := unpackDomainName(data, int(lablePtr), maxDepth-1)
            if tempErr != nil {
                return
            }
            name += tempName
            return
        default:
            // 0x80 and 0x40 are reserved
            err = NewError("fomart error")
            return
        }
    }
    return
}

func unpackQuestion(data []byte, index int, question *Question) (offset int, err error) {
    offset = index
    question.Name, offset, err = unpackDomainName(data, offset, MAX_COMPRESSION_DEPTH)
    if err != nil {
        return
    }
    if offset+4 > len(data) {
        err = NewError("out of range")
        return
    }
    question.Type, offset = unpackUint16(data, offset)
    question.Class, offset = unpackUint16(data, offset)
    return
}

func unpackRR(data []byte, index int) (rr RR, offset int, err error) {
    offset = index
    var hdr RRHeader
    hdr.Name, offset, err = unpackDomainName(data, offset, MAX_COMPRESSION_DEPTH)
    if err != nil {
        return
    }
    if offset+10 > len(data) {
        err = NewError("out of range")
        return
    }
    hdr.Type, offset = unpackUint16(data, offset)
    hdr.Class, offset = unpackUint16(data, offset)
    hdr.Ttl, offset = unpackUint32(data, offset)
    hdr.RDLength, offset = unpackUint16(data, offset)

    if hdr.Class != CLASS_INET {
        err = NewError("unimplement")
        return
    }

    rr, err = RRNew(hdr.Type)
    if err != nil {
        return
    }
    rr.SetHeader(&hdr)
    offset, err = rr.UnpackRData(data, offset)
    return
}

func packUint16(number uint16, data []byte, index int) (offset int) {
    data[index] = uint8(number >> 8)
    data[index+1] = uint8(number)
    offset = index + 2
    return
}
func packUint32(number uint32, data []byte, index int) (offset int) {
    data[index] = uint8(number >> 24)
    data[index+1] = uint8(number >> 16)
    data[index+2] = uint8(number >> 8)
    data[index+3] = uint8(number)
    offset = index + 4
    return
}

func (message *Message) Pack(data []byte, needCompress bool) (length int, err error) {
    var compression map[string]int
    if needCompress {
        compression = make(map[string]int) // Compression pointer mappings
    } else {
        compression = nil
    }

    length = 0
    dataLen := len(data)
    if dataLen < HEADER_LENGTH {
        err = NewError("too short")
        return
    }
    length = packUint16(message.Hdr.Id, data, length)
    if message.Hdr.QueryResponse == QR_RESPONSE {
        data[length] |= _QUREY_RESPONSE
    }
    data[length] |= message.Hdr.Opcode << 4 >> 1
    if message.Hdr.AuthAnswer {
        data[length] |= _AUTH_ANSWER
    }
    if message.Hdr.Truncated {
        data[length] |= _TRUNCATED
    }
    if message.Hdr.RecursionDesired {
        data[length] |= _RECURSION_DESIRED
    }
    length++
    if message.Hdr.RecursionAvailable {
        data[length] |= _RECURSION_AVAILABLE
    }
    data[length] |= message.Hdr.Zero << 5 >> 1
    data[length] |= message.Hdr.Rcode & _RCODE
    length++

    message.Hdr.QDCount = uint16(len(message.Question))
    message.Hdr.ANCount = uint16(len(message.Answer))
    message.Hdr.NSCount = uint16(len(message.Authority))
    message.Hdr.ARCount = uint16(len(message.Additional))
    length = packUint16(message.Hdr.QDCount, data, length)
    length = packUint16(message.Hdr.ANCount, data, length)
    length = packUint16(message.Hdr.NSCount, data, length)
    length = packUint16(message.Hdr.ARCount, data, length)
    for _, question := range message.Question {
        length, err = packQuestion(data, length, &question, compression)
        if err != nil {
            return
        }
    }
    length, err = packRRs(message.Answer, data, length, compression)
    if err != nil {
        return
    }
    //length, err = packRRs(message.Authority, data, length)
    //if err != nil {
    //return
    //}
    //length, err = packRRs(message.Additional, data, length)
    //if err != nil {
    //return
    //}
    return
}

func packQuestion(buf []byte, index int, question *Question, compression map[string]int) (offset int, err error) {
    offset, err = packDomainName(question.Name, buf, index, compression)
    if err != nil {
        return
    }
    if 4 > len(buf)-offset {
        err = NewError("buffer too small to store question type and class")
        return
    }
    offset = packUint16(question.Type, buf, offset)
    offset = packUint16(question.Class, buf, offset)
    return
}

func packDomainName(name string, buf []byte, index int, compression map[string]int) (offset int, err error) {
    bufLen := len(buf)
    offset = index
    if len(name) > MAX_DOMAIN_NAME_LEN {
        err = NewError(fmt.Sprintf("Domain name length must <= %d: %s", MAX_DOMAIN_NAME_LEN,
            name))
        return
    }
    tempName := name
    for {
        if len(tempName) == 0 {
            break
        }
        if compression != nil {
            //need compress
            if ptr, ok := compression[tempName]; ok {
                offset = packUint16(uint16(ptr), buf, offset)
                buf[offset-2] |= 0xC0
                return
            }
        }
        dotIndex := strings.Index(tempName, ".")
        var label string
        if dotIndex == -1 {
            label = tempName
            tempName = ""
        } else {
            label = tempName[:dotIndex]
            tempName = tempName[dotIndex+1:]
        }
        labelLen := len(label)
        if labelLen > MAX_DOMAIN_LABEL_LEN {
            err = NewError(fmt.Sprintf("Domain label length must <= %d: %s", MAX_DOMAIN_LABEL_LEN, label))
            return
        }
        if labelLen+1 > bufLen-offset {
            err = NewError("buffer too small to store " + name)
            return
        }
        if compression != nil {
            compression[label+"."+tempName] = offset
        }
        buf[offset] = uint8(labelLen)
        offset++
        copy(buf[offset:], label)
        offset += labelLen
    }
    //labels := strings.Split(name, ".")
    //for _, label := range labels {
    //labelLen := len(label)
    //if labelLen > MAX_DOMAIN_LABEL_LEN {
    //err = NewError(fmt.Sprintf("Domain label length must <= %d: %s", MAX_DOMAIN_LABEL_LEN, label))
    //return
    //}
    //if labelLen+1 > bufLen-offset {
    //err = NewError("buffer too small to store " + name)
    //return
    //}
    //buf[offset] = uint8(labelLen)
    //copy(buf[offset+1:], label)
    //offset += 1 + labelLen
    //}
    buf[offset] = 0
    offset++
    return
}

func packRRs(rrs []RR, buf []byte, index int, compression map[string]int) (offset int, err error) {
    offset = index
    for _, rr := range rrs {
        offset, err = packRR(rr, buf, offset, compression)
        if err != nil {
            return
        }
    }
    return
}
func packRR(rr RR, buf []byte, index int, compression map[string]int) (offset int, err error) {
    offset = index
    header := rr.Header()
    offset, err = packDomainName(header.Name, buf, offset, compression)
    if err != nil {
        return
    }
    //10 = sizeof(Type)+sizeof(Class)+sizeof(Ttl)+sizeof(RDLength)
    if 10 > len(buf)-offset {
        err = NewError("buffer too small to store RR")
        return
    }
    offset = packUint16(header.Type, buf, offset)
    offset = packUint16(header.Class, buf, offset)
    offset = packUint32(header.Ttl, buf, offset)
    start := offset + 2
    offset, err = rr.PackRData(buf, start)
    if err != nil {
        return
    }
    header.RDLength = uint16(offset - start)
    packUint16(header.RDLength, buf, start-2)
    return
}
