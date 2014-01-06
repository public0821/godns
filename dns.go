package main

import (
	//"./dns"
	"fmt"
	"net"
    "os"
    "bufio"
    "strings"
)

func getDnsServer()(servers []string, err error){
    file, err := os.Open("/etc/resolv.conf")
    if err != nil {
        return
    }
    defer file.Close()
    reader := bufio.NewReader(file)
    lineBytes,_, err := reader.ReadLine()
    for err == nil {
        fields := strings.Fields(string(lineBytes))
        if len(fields) >= 2 && fields[0] == "nameserver" {
            fmt.Println(fields)
            servers = append(servers, fields[1])
        }
        lineBytes, _, err = reader.ReadLine()
    }
    return
}

func main() {
    fmt.Println(getDnsServer())
	addr := net.UDPAddr{Port: 5354}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := make([]byte, 65535)

	for {
		length, client, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
			conn.Close()
			return
		}
		fmt.Println(client)
		fmt.Println(string(buf[:length]))
	}
	fmt.Println(addr)
}
