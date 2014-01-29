package util

import (
    //"./dns"
    "bufio"
    "io"
    "os"
    "strings"
)

func GetDnsServer() (servers []string, err error) {
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
