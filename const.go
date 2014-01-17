package main

const (
    MAX_DOMAIN_NAME_LEN  = 255
    MAX_DOMAIN_LABEL_LEN = 63
    MAX_UDP_MESSAGE_LEN  = 512
)

const (
    QR_REQUEST  uint8 = 0
    QR_RESPONSE uint8 = 1

    // valid RRHeader.TYPE_ and Question.TYPE_
    TYPE_None       uint16 = 0
    TYPE_A          uint16 = 1
    TYPE_NS         uint16 = 2
    TYPE_MD         uint16 = 3
    TYPE_MF         uint16 = 4
    TYPE_CNAME      uint16 = 5
    TYPE_SOA        uint16 = 6
    TYPE_MB         uint16 = 7
    TYPE_MG         uint16 = 8
    TYPE_MR         uint16 = 9
    TYPE_NULL       uint16 = 10
    TYPE_WKS        uint16 = 11
    TYPE_PTR        uint16 = 12
    TYPE_HINFO      uint16 = 13
    TYPE_MINFO      uint16 = 14
    TYPE_MX         uint16 = 15
    TYPE_TXT        uint16 = 16
    TYPE_RP         uint16 = 17
    TYPE_AFSDB      uint16 = 18
    TYPE_X25        uint16 = 19
    TYPE_ISDN       uint16 = 20
    TYPE_RT         uint16 = 21
    TYPE_NSAP       uint16 = 22
    TYPE_NSAPPTR    uint16 = 23
    TYPE_SIG        uint16 = 24
    TYPE_KEY        uint16 = 25
    TYPE_PX         uint16 = 26
    TYPE_GPOS       uint16 = 27
    TYPE_AAAA       uint16 = 28
    TYPE_LOC        uint16 = 29
    TYPE_NXT        uint16 = 30
    TYPE_EID        uint16 = 31
    TYPE_NIMLOC     uint16 = 32
    TYPE_SRV        uint16 = 33
    TYPE_ATMA       uint16 = 34
    TYPE_NAPTR      uint16 = 35
    TYPE_KX         uint16 = 36
    TYPE_CERT       uint16 = 37
    TYPE_DNAME      uint16 = 39
    TYPE_OPT        uint16 = 41 // EDNS
    TYPE_DS         uint16 = 43
    TYPE_SSHFP      uint16 = 44
    TYPE_IPSECKEY   uint16 = 45
    TYPE_RRSIG      uint16 = 46
    TYPE_NSEC       uint16 = 47
    TYPE_DNSKEY     uint16 = 48
    TYPE_DHCID      uint16 = 49
    TYPE_NSEC3      uint16 = 50
    TYPE_NSEC3PARAM uint16 = 51
    TYPE_TLSA       uint16 = 52
    TYPE_HIP        uint16 = 55
    TYPE_NINFO      uint16 = 56
    TYPE_RKEY       uint16 = 57
    TYPE_TALINK     uint16 = 58
    TYPE_CDS        uint16 = 59
    TYPE_SPF        uint16 = 99
    TYPE_UINFO      uint16 = 100
    TYPE_UID        uint16 = 101
    TYPE_GID        uint16 = 102
    TYPE_UNSPEC     uint16 = 103
    TYPE_NID        uint16 = 104
    TYPE_L32        uint16 = 105
    TYPE_L64        uint16 = 106
    TYPE_LP         uint16 = 107
    TYPE_EUI48      uint16 = 108
    TYPE_EUI64      uint16 = 109

    TYPE_TKEY uint16 = 249
    TYPE_TSIG uint16 = 250

    // valid Question.TYPE_ only
    TYPE_IXFR  uint16 = 251
    TYPE_AXFR  uint16 = 252
    TYPE_MAILB uint16 = 253
    TYPE_MAILA uint16 = 254
    TYPE_ANY   uint16 = 255

    TYPE_URI      uint16 = 256
    TYPE_CAA      uint16 = 257
    TYPE_TA       uint16 = 32768
    TYPE_DLV      uint16 = 32769
    TYPE_Reserved uint16 = 65535

    // valid Question.CLASS_ and RRHeader.CLASS_
    CLASS_INET   uint16 = 1
    CLASS_CSNET  uint16 = 2
    CLASS_CHAOS  uint16 = 3
    CLASS_HESIOD uint16 = 4
    CLASS_NONE   uint16 = 254
    CLASS_ANY    uint16 = 255

    // Message.Hdr.Rcode
    RCODE_SUCCESS        uint8 = 0
    RCODE_FORMATERROR    uint8 = 1
    RCODE_SERVERFAILURE  uint8 = 2
    RCODE_NAMEERROR      uint8 = 3
    RCODE_NOTIMPLEMENTED uint8 = 4
    RCODE_REFUSED        uint8 = 5
    RCODE_YXDOMAIN       uint8 = 6
    RCODE_YXRRSET        uint8 = 7
    RCODE_NXRRSET        uint8 = 8
    RCODE_NOTAUTH        uint8 = 9
    RCODE_NOTZONE        uint8 = 10
    RCODE_BADSIG         uint8 = 16 // TSIG
    RCODE_BADVERS        uint8 = 16 // EDNS0
    RCODE_BADKEY         uint8 = 17
    RCODE_BADTIME        uint8 = 18
    RCODE_BADMODE        uint8 = 19 // TKEY
    RCODE_BADNAME        uint8 = 20
    RCODE_BADALG         uint8 = 21
    RCODE_BADTRUNC       uint8 = 22 // TSIG

    // OPCODE_
    OPCODE_QUERY  uint8 = 0
    OPCODE_IQUERY uint8 = 1
    OPCODE_STATUS uint8 = 2
    // There is no 3
    OPCODE_NOTIFY uint8 = 4
    OPCODE_UPDATE uint8 = 5
)

const (
    _QUREY_RESPONSE      uint8 = 0x80
    _OPCODE              uint8 = 0x78
    _AUTH_ANSWER         uint8 = 0x04
    _TRUNCATED           uint8 = 0x02
    _RECURSION_DESIRED   uint8 = 0x01
    _RECURSION_AVAILABLE uint8 = 0x80
    _ZERO                uint8 = 0x70
    _RCODE               uint8 = 0x0f
)

// Map of strings for each RR wire type.
var TypeToString = map[uint16]string{
    TYPE_A:          "A",
    TYPE_AAAA:       "AAAA",
    TYPE_AFSDB:      "AFSDB",
    TYPE_ANY:        "ANY", // Meta RR
    TYPE_ATMA:       "ATMA",
    TYPE_AXFR:       "AXFR", // Meta RR
    TYPE_CAA:        "CAA",
    TYPE_CDS:        "CDS",
    TYPE_CERT:       "CERT",
    TYPE_CNAME:      "CNAME",
    TYPE_DHCID:      "DHCID",
    TYPE_DLV:        "DLV",
    TYPE_DNAME:      "DNAME",
    TYPE_DNSKEY:     "DNSKEY",
    TYPE_DS:         "DS",
    TYPE_EID:        "EID",
    TYPE_EUI48:      "EUI48",
    TYPE_EUI64:      "EUI64",
    TYPE_GID:        "GID",
    TYPE_GPOS:       "GPOS",
    TYPE_HINFO:      "HINFO",
    TYPE_HIP:        "HIP",
    TYPE_IPSECKEY:   "IPSECKEY",
    TYPE_ISDN:       "ISDN",
    TYPE_IXFR:       "IXFR", // Meta RR
    TYPE_KX:         "KX",
    TYPE_L32:        "L32",
    TYPE_L64:        "L64",
    TYPE_LOC:        "LOC",
    TYPE_LP:         "LP",
    TYPE_MB:         "MB",
    TYPE_MD:         "MD",
    TYPE_MF:         "MF",
    TYPE_MG:         "MG",
    TYPE_MINFO:      "MINFO",
    TYPE_MR:         "MR",
    TYPE_MX:         "MX",
    TYPE_NAPTR:      "NAPTR",
    TYPE_NID:        "NID",
    TYPE_NINFO:      "NINFO",
    TYPE_NIMLOC:     "NIMLOC",
    TYPE_NS:         "NS",
    TYPE_NSAP:       "NSAP",
    TYPE_NSAPPTR:    "NSAP-PTR",
    TYPE_NSEC3:      "NSEC3",
    TYPE_NSEC3PARAM: "NSEC3PARAM",
    TYPE_NSEC:       "NSEC",
    TYPE_NULL:       "NULL",
    TYPE_OPT:        "OPT",
    TYPE_PTR:        "PTR",
    TYPE_RKEY:       "RKEY",
    TYPE_RP:         "RP",
    TYPE_RRSIG:      "RRSIG",
    TYPE_RT:         "RT",
    TYPE_SOA:        "SOA",
    TYPE_SPF:        "SPF",
    TYPE_SRV:        "SRV",
    TYPE_SSHFP:      "SSHFP",
    TYPE_TA:         "TA",
    TYPE_TALINK:     "TALINK",
    TYPE_TKEY:       "TKEY", // Meta RR
    TYPE_TLSA:       "TLSA",
    TYPE_TSIG:       "TSIG", // Meta RR
    TYPE_TXT:        "TXT",
    TYPE_PX:         "PX",
    TYPE_UID:        "UID",
    TYPE_UINFO:      "UINFO",
    TYPE_UNSPEC:     "UNSPEC",
    TYPE_URI:        "URI",
    TYPE_WKS:        "WKS",
    TYPE_X25:        "X25",
}

// Map of strings for each CLASS wire type.
var ClassToString = map[uint16]string{
    CLASS_INET:   "IN",
    CLASS_CSNET:  "CS",
    CLASS_CHAOS:  "CH",
    CLASS_HESIOD: "HS",
    CLASS_NONE:   "NONE",
    CLASS_ANY:    "ANY",
}

// Map of strings for opcodes.
var OpcodeToString = map[uint8]string{
    OPCODE_QUERY:  "QUERY",
    OPCODE_IQUERY: "IQUERY",
    OPCODE_STATUS: "STATUS",
    OPCODE_NOTIFY: "NOTIFY",
    OPCODE_UPDATE: "UPDATE",
}

// Map of strings for rcodes.
var RcodeToString = map[uint8]string{
    RCODE_SUCCESS:        "NOERROR",
    RCODE_FORMATERROR:    "FORMERR",
    RCODE_SERVERFAILURE:  "SERVFAIL",
    RCODE_NAMEERROR:      "NXDOMAIN",
    RCODE_NOTIMPLEMENTED: "NOTIMPL",
    RCODE_REFUSED:        "REFUSED",
    RCODE_YXDOMAIN:       "YXDOMAIN", // From RFC 2136
    RCODE_YXRRSET:        "YXRRSET",
    RCODE_NXRRSET:        "NXRRSET",
    RCODE_NOTAUTH:        "NOTAUTH",
    RCODE_NOTZONE:        "NOTZONE",
    RCODE_BADSIG:         "BADSIG", // Also known as RcodeBadVers, see RFC 6891
    //     RCODE_BadVers:        "BADVERS",
    RCODE_BADKEY:   "BADKEY",
    RCODE_BADTIME:  "BADTIME",
    RCODE_BADMODE:  "BADMODE",
    RCODE_BADNAME:  "BADNAME",
    RCODE_BADALG:   "BADALG",
    RCODE_BADTRUNC: "BADTRUNC",
}

// Reverse a map
func reverseInt8(m map[uint8]string) map[string]uint8 {
    n := make(map[string]uint8)
    for k, v := range m {
        n[v] = k
    }
    return n
}

func reverseInt16(m map[uint16]string) map[string]uint16 {
    n := make(map[string]uint16)
    for k, v := range m {
        n[v] = k
    }
    return n
}

// Reverse, needed for string parsing.
var StringToTYPE_ = reverseInt16(TypeToString)
var StringToCLASS_ = reverseInt16(ClassToString)

// Map of opcodes strings.
var StringToOpcode = reverseInt8(OpcodeToString)

// Map of rcodes strings.
var StringToRcode = reverseInt8(RcodeToString)
