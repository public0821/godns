package main

import (
	//"./dns"
    "io"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

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
        commentIndex := strings.Index(line, "#")
        if commentIndex != -1 {
            line = line[:commentIndex]
        }
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "nameserver" {
			servers = append(servers, fields[1])
		}
		lineBytes, _, err = reader.ReadLine()
	}
    if err == io.EOF {
        err = nil
    }
	return
}

func main() {
    dnsServers, err := getDnsServer()
    fmt.Println(dnsServers, err)
    if err != nil {
        fmt.Println(err)
        return
    }
	addr := net.UDPAddr{Port: 53}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println(err)
		return
	}
    defer conn.Close()
	buf := make([]byte, 65535)

    //var remote net.UDPConn
    localAddr := net.UDPAddr{Port:12345}
    remote, err := net.ListenUDP("udp", &localAddr)
    if err != nil {
        fmt.Println(err)
        return
    }
    var remoteAddr net.UDPAddr
    remoteAddr.IP = net.ParseIP(dnsServers[0])
    remoteAddr.Port = 53
	for {
		length, client, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
            continue
		}
        _,err = remote.WriteToUDP(buf[:length], &remoteAddr) 
		if err != nil {
			fmt.Println(err)
            continue
		}
        length, _, err = conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
            continue
		}
        _, err = conn.WriteToUDP(buf[:length], client)
		if err != nil {
			fmt.Println(err)
            continue
		}
	}
	fmt.Println(addr)
}
