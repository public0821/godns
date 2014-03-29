package db

import (
//"github.com/public0821/dnserver/errors"
)

type Record struct {
    Id    int64
    Name  string
    Class uint16
    Type  uint16
    Value string
    Ttl   uint32
}

type User struct {
    Name string
    Pwd  string
}

type SysOption struct {
    Name  string
    Value string
}

type ForwardServer struct {
    Ip string
}
