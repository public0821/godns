package dns

type RRType int
type OpcodeType int
type RcodeType int

const (
        RR_REQUEST  RRType = 0
        RR_RESPONSE RRType = 1
)

type Header struct {
        Id                 uint16
        RRCode             RRType //request or response
        Opcode             OpcodeType
        AuthAnswer         bool
        Truncated          bool
        RecursionDesired   bool
        RecursionAvailable bool
        Rcode              RcodeType
}

type Question struct {
        Name   string `dns:"cdomain-name"` // "cdomain-name" specifies encoding (and may be compressed)
        Qtype  uint16
        Qclass uint16
}

type RRBase struct {
        Name     string `dns:"cdomain-name"`
        Type     uint16
        Class    uint16
        Ttl      uint32
        Rdlength uint16 // length of data after header
}

type A struct {
        RRBase
        IPv4 net.IP
}

type Message struct {
        Hdr        Header
        Question   []Question // Holds the RR(s) of the question section.
        Answer     []RRBase      // Holds the RR(s) of the answer section.
        Authority  []RRBase       // Holds the RR(s) of the authority section.
        Additional []RRBase       // Holds the RR(s) of the additional section.
}
