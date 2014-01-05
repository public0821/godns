package main

import (
	"./dns"
	"fmt"
	"net"
)

func main() {
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
