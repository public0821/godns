package db

import (
//"github.com/public0821/dnserver/errors"
)

type RRecord struct {
    Id    int64
    Name  string
    Class uint16
    Type  uint16
    Value string
    Ttl   uint32
}

type User struct {
    Id   int
    Name string
    Pwd  string
}

type SysOption struct {
    Id    int
    Name  string
    Value string
}

type ForwardServer struct {
    Id  int
    Ip  string
}
