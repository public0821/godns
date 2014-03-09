package dns

import (
	"fmt"
	"github.com/public0821/dnserver/errors"
	"strconv"
	"strings"
)

type RRMx struct {
	rrBase
	priority uint16
	domain   string
}

func (m *RRMx) PackRData(buf []byte, index int) (offset int, err error) {
	offset = index
	if len(buf) < offset+2 {
		err = errors.New("buf too short")
		return
	}
	offset = PackUint16(m.priority, buf, offset)
	offset, err = PackDomainName(m.domain, buf, offset, nil)
	if err != nil {
		return
	}
	return
}

func (m *RRMx) UnpackRData(buf []byte, index int) (offset int, err error) {
	err = errors.New("unsupport unpack mx record ")
	return
}

func (m *RRMx) SetRData(data string) (err error) {
	fields := strings.Split(data, " ")
	if len(fields) != 2 {
		err = errors.New("invalid mx rdata: " + data)
		return
	}
	m.domain = fields[0]
	priority, err := strconv.Atoi(fields[1])
	if err != nil {
		err = errors.New("invalid mx priority: " + fields[1])
		return
	}
	m.priority = uint16(priority)
	return
}

func (m *RRMx) String() (text string) {
	text += m.Hdr.String()
	text += fmt.Sprintf("Domain:\t%s\n", m.domain)
	text += fmt.Sprintf("Priority:\t%d\n", m.priority)
	return
}
