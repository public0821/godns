package main

import (
	//"./dns"
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const FORWARD_PORT_NUM = 5
const MAX_BUF_LEN = 65535

func getDnsServer() (servers []string, err error) {
	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	lineBytes, _, err := reader.ReadLine()
	for err == nil {
		line := string(lineBytes)
		//remove comments
		commentIndex := strings.Index(line, "#")
		if commentIndex != -1 {
			line = line[:commentIndex]
		}

		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == "nameserver" {
			servers = append(servers, fields[1])
		}
		lineBytes, _, err = reader.ReadLine()
	}
	if err == io.EOF {
		err = nil
	}
	return
}

var session Session

func doForwardToResolver(forwardConns []*net.UDPConn, resolverAddrs []string, recvMsgChannel chan RecvMsg) {
	log.Println("start doForwardToResolver")
	addrIndex := 0
	connIndex := 0
	for recvMsg := range recvMsgChannel {
		addr, err := net.ResolveUDPAddr("udp", resolverAddrs[addrIndex])
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("write request to resolver")
		conn := forwardConns[connIndex]
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
		key.name = "test"
		key.id = 1

		var value SessionValue
		value.clientAddr = recvMsg.addr
		value.id = 2
		value.name = "test"

		session.Lock()
		session.buffer[key] = value
		session.Unlock()
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
		var key SessionKey
		key.forwardConn = conn
		key.name = "test"
		key.id = 1

		session.Lock()
		value, ok := session.buffer[key]
		if ok {
			delete(session.buffer, key)
		}
		session.Unlock()

		if ok {
			log.Println("write response to client")
			_, err = server.WriteToUDP(buf[:length], &value.clientAddr)
			if err != nil {
				log.Println(err)
				continue
			}
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
	dnsServers, err := getDnsServer()
	if err != nil {
		log.Panicln(err)
		return
	}
	var resolverAddrs []string
	for _, dnsServer := range dnsServers {
		resolverAddrs = append(resolverAddrs, dnsServer+":53")
	}

	//initialize session
	session.buffer = make(map[SessionKey]SessionValue)

	//listen for service
	addr := net.UDPAddr{Port: 53}
	server, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Panicln(err)
		return
	}
	defer server.Close()

	//initialize forward conn
	forwardConns := make([]*net.UDPConn, FORWARD_PORT_NUM)
	for i := 0; i < FORWARD_PORT_NUM; i++ {
		forwardConn, err := net.ListenUDP("udp", &net.UDPAddr{})
		if err != nil {
			log.Panicln(err)
			return
		}
		forwardConns[i] = forwardConn
		go doRecvFromResolver(server, forwardConn)
		defer forwardConn.Close()
	}

	recvMsgChannel := make(chan RecvMsg, 10)
	go doForwardToResolver(forwardConns, resolverAddrs, recvMsgChannel)

	buf := make([]byte, 65535)
	for {
		length, client, err := server.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		recvMsgChannel <- RecvMsg{addr: *client, data: buf[:length]}
	}
}
