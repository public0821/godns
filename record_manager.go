package main

import (
//"fmt"
)

type Record struct {
    Id    int64
    Name  string
    Class uint16
    Type  uint16
    Value string
    Ttl   uint32
}

type RecordManager interface {
    Open() (err error)
    DeleteAllRecord() (err error)
    AddRecord(record *Record) (err error)
    DeleteRecord(id int64) (err error)
    ModifyRecord(record *Record) (err error)
    QueryRecord(record *Record) (records []Record, err error)
    Count() (count int, err error)
    Close() (err error)
}
