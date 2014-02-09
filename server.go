package main

import (
	//"./dns"
	//"bufio"
	"github.com/public0821/dnserver/db"
	"github.com/public0821/dnserver/dns"
	"github.com/public0821/dnserver/util"
	"github.com/public0821/dnserver/web"
	//"io"
	"log"
	"math/rand"
	"net"
	"os"
	//"strings"
	"time"
)

const FORWARD_PORT_NUM = 5
const MAX_BUF_LEN = 65535
const SESSION_TIMEOUT = 10

var gsession Session

func RRConstruct(rrecord *db.RRecord) (rr dns.RR, err error) {
	rr, err = dns.RRNew(rrecord.Type)
	if err != nil {
		return
	}
	var header dns.RRHeader
	header.Name = rrecord.Name
	header.Class = rrecord.Class
	header.Type = rrecord.Type
	header.Ttl = rrecord.Ttl
	rr.SetHeader(&header)
	err = rr.SetRData(rrecord.Value)
	return
}

func doCleanTimeoutSession() {
	for {
		now := time.Now().Unix()
		gsession.Lock()
		for key, value := range gsession.buffer {
			if now-value.time >= SESSION_TIMEOUT {
				delete(gsession.buffer, key)
			}
		}
		gsession.Unlock()
		time.Sleep(time.Second * (SESSION_TIMEOUT / 2))
	}
}

func doForwardToResolver(server *net.UDPConn, forwardConns []*net.UDPConn, resolverAddrs []string, recvMsgChannel chan RecvMsg) {
	buf := make([]byte, MAX_BUF_LEN)
	log.Println("start doForwardToResolver")
	addrIndex := 0
	connIndex := 0
	rand.Seed(time.Now().UTC().UnixNano())
	for recvMsg := range recvMsgChannel {
		//unpack message
		var msg dns.Message
		_, err := msg.UnpackHeaderAndQuestion(recvMsg.data)
		if err != nil {
			log.Println(err)
			continue
		}
		//check whether the domain name is in record manager
		question := &msg.Question[0]
		var rrecord db.RRecord
		rrecord.Name = question.Name
		rrecord.Class = question.Class
		rrecord.Type = question.Type
		records, err := db.Query(&rrecord, 0, 0)
		if err != nil {
			log.Println(err)
			continue
		}
		//construct an answer record and send to client
		if len(records) > 0 {
			tempRRecord, _ := records[0].(db.RRecord)
			rr, err := RRConstruct(&tempRRecord)
			if err != nil {
				log.Println(err)
				continue
			}
			msg.Hdr.Rcode = dns.RCODE_SUCCESS
			msg.Hdr.QueryResponse = dns.QR_RESPONSE
			msg.Answer = append(msg.Answer, rr)
			buflen, err := msg.Pack(buf, true)
			log.Println("construct an answer record")
			_, err = server.WriteToUDP(buf[:buflen], &recvMsg.addr)
			if err != nil {
				log.Println(err)
				continue
			}
			continue
		}
		//forward to upstream resolver
		addr, err := net.ResolveUDPAddr("udp", resolverAddrs[addrIndex])
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("write request to resolver")
		conn := forwardConns[connIndex]
		//modify id
		newId := uint16(rand.Int())
		recvMsg.data[0] = uint8(newId >> 8)
		recvMsg.data[1] = uint8(newId)
		_, err = conn.WriteToUDP(recvMsg.data, addr)
		if err != nil {
			log.Println(err)
			continue
		}
		connIndex += 1
		if connIndex >= len(forwardConns) {
			connIndex = 0
		}
		addrIndex += 1
		if addrIndex >= len(resolverAddrs) {
			addrIndex = 0
		}

		var key SessionKey
		key.forwardConn = conn
		key.name = question.Name
		key.id = newId

		var value SessionValue
		value.clientAddr = recvMsg.addr
		value.id = msg.Hdr.Id
		value.time = time.Now().Unix()

		gsession.Lock()
		gsession.buffer[key] = value
		gsession.Unlock()
	}
}

func doRecvFromResolver(server, conn *net.UDPConn) {
	log.Println("start  doRecvFromResolver")
	buf := make([]byte, MAX_BUF_LEN)
	for {
		length, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("receive response from resolver")
		//unpack message
		var msg dns.Message
		_, err = msg.UnpackHeaderAndQuestion(buf[:length])
		if err != nil {
			log.Println(err)
			continue
		}

		var key SessionKey
		key.forwardConn = conn
		question := &msg.Question[0]
		key.name = question.Name
		key.id = msg.Hdr.Id

		gsession.Lock()
		value, ok := gsession.buffer[key]
		if ok {
			delete(gsession.buffer, key)
			gsession.Unlock()
		} else {
			log.Println("key not find in session: ", key)
			gsession.Unlock()
			continue
		}

		log.Println("write response to client")
		//modify id
		newId := value.id
		buf[0] = uint8(newId >> 8)
		buf[1] = uint8(newId)
		_, err = server.WriteToUDP(buf[:length], &value.clientAddr)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

type RecvMsg struct {
	addr net.UDPAddr
	data []byte
}

func main() {
	//initialize log
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//initialize resolver addrs
	dnsServers, err := util.GetDnsServer()
	if err != nil {
		log.Println(err)
		return
	}
	var resolverAddrs []string
	for _, dnsServer := range dnsServers {
		resolverAddrs = append(resolverAddrs, dnsServer+":53")
	}

	//initialize session
	gsession.buffer = make(map[SessionKey]SessionValue)

	//initialize db
	_, err = db.NewDBManager()
	if err != nil {
		log.Println(err)
		return
	}

	//listen for service
	addr := net.UDPAddr{Port: 53}
	server, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Println(err)
		return
	}
	defer server.Close()

	//initialize forward conn
	forwardConns := make([]*net.UDPConn, FORWARD_PORT_NUM)
	for i := 0; i < FORWARD_PORT_NUM; i++ {
		forwardConn, err := net.ListenUDP("udp", &net.UDPAddr{})
		if err != nil {
			log.Println(err)
			return
		}
		forwardConns[i] = forwardConn
		go doRecvFromResolver(server, forwardConn)
		defer forwardConn.Close()
	}

	recvMsgChannel := make(chan RecvMsg, 10)
	go doForwardToResolver(server, forwardConns, resolverAddrs, recvMsgChannel)

	//start web
	go web.Start()

	buf := make([]byte, 65535)
	for {
		length, client, err := server.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		//log.Println(buf[:length])
		recvMsgChannel <- RecvMsg{addr: *client, data: buf[:length]}
	}
}
