package main

const (
	MAX_DOMAIN_NAME_LEN  = 255
	MAX_DOMAIN_LABEL_LEN = 63
	MAX_UDP_MESSAGE_LEN  = 512
)

const (
	QR_REQUEST  uint8 = 0
	QR_RESPONSE uint8 = 1

	// valid RRHeader.Type and Question.Type
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

	// valid Question.Type only
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

	// valid Question.Class and RRHeader.Class
	CLASS_INET   = 1
	CLASS_CSNET  = 2
	CLASS_CHAOS  = 3
	CLASS_HESIOD = 4
	CLASS_NONE   = 254
	CLASS_ANY    = 255

	// Message.Hdr.Rcode
	RCODE_SUCCESS        = 0
	RCODE_FORMATERROR    = 1
	RCODE_SERVERFAILURE  = 2
	RCODE_NAMEERROR      = 3
	RCODE_NOTIMPLEMENTED = 4
	RCODE_REFUSED        = 5
	RCODE_YXDOMAIN       = 6
	RCODE_YXRRSET        = 7
	RCODE_NXRRSET        = 8
	RCODE_NOTAUTH        = 9
	RCODE_NOTZONE        = 10
	RCODE_BADSIG         = 16 // TSIG
	RCODE_BADVERS        = 16 // EDNS0
	RCODE_BADKEY         = 17
	RCODE_BADTIME        = 18
	RCODE_BADMODE        = 19 // TKEY
	RCODE_BADNAME        = 20
	RCODE_BADALG         = 21
	RCODE_BADTRUNC       = 22 // TSIG

	// OPCODE_
	OPCODE_QUERY  = 0
	OPCODE_IQUERY = 1
	OPCODE_STATUS = 2
	// There is no 3
	OPCODE_NOTIFY = 4
	OPCODE_UPDATE = 5
)

const (
	_QUREY_RESPONSE      = 0x80
	_OPCODE              = 0x78
	_AUTH_ANSWER         = 0x04
	_TRUNCATED           = 0x02
	_RECURSION_DESIRED   = 0x01
	_RECURSION_AVAILABLE = 0x80
	_ZERO                = 0x70
	_RCODE               = 0x0f
)
