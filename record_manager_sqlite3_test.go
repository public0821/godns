package main

import (
    //"fmt"
    "testing"
)

func TestAdd(t *testing.T) {
    var m RecordManager
    m = &RecordManagerSqlite3{}
    err := m.Open()
    if err != nil {
        t.Error(err)
        return
    }
    defer m.Close()

    count, err := m.Count()
    if err != nil {
        t.Error(err)
        return
    }

    var record Record
    record.Name = "www.test.com"
    record.Class = 1
    record.Type = 1
    record.Value = "10.32.171.60"
    err = m.AddRecord(&record)
    if err != nil {
        t.Error(err)
        return
    }
    newCount, err := m.Count()
    if err != nil {
        t.Error(err)
        return
    }
    if newCount-count != 1 {
        t.Error("add failed")
        return
    }
}

func TestQuery(t *testing.T) {
    var m RecordManager
    m = &RecordManagerSqlite3{}
    err := m.Open()
    if err != nil {
        t.Error(err)
        return
    }
    defer m.Close()
    var record Record
    record.Name = "www.test.com"
    record.Class = 1
    record.Type = 1
    record.Value = "10.32.171.60"
    records, err := m.QueryRecord(&record)
    if err != nil {
        t.Error(err)
        return
    }
    if len(records) != 1 {
        t.Error("query failed")
        return
    }
}

func TestModify(t *testing.T) {
    var m RecordManager
    m = &RecordManagerSqlite3{}
    err := m.Open()
    if err != nil {
        t.Error(err)
        return
    }
    defer m.Close()
    var record Record
    record.Name = "www.test.com"
    record.Class = 1
    record.Type = 1
    record.Value = "10.32.171.60"
    records, err := m.QueryRecord(&record)
    if err != nil {
        t.Error(err)
        return
    }
    if len(records) != 1 {
        t.Error("query failed")
        return
    }
    records[0].Name = "www.test1.com"
    records[0].Class = 2
    records[0].Type = 2
    records[0].Value = "11.32.171.60"
    err = m.ModifyRecord(&records[0])
    if err != nil {
        t.Error(err)
        return
    }
    newRecords, err := m.QueryRecord(&records[0])
    if err != nil {
        t.Error(err)
        return
    }
    if len(newRecords) != 1 {
        t.Error("modify record failed")
        return
    }

}

func TestDelete(t *testing.T) {
    TestAdd(t)
    var m RecordManager
    m = &RecordManagerSqlite3{}
    err := m.Open()
    if err != nil {
        t.Error(err)
        return
    }
    defer m.Close()
    count, err := m.Count()
    if err != nil {
        t.Error(err)
        return
    }
    var record Record
    record.Name = "www.test.com"
    record.Class = 1
    record.Type = 1
    record.Value = "10.32.171.60"
    records, err := m.QueryRecord(&record)
    if err != nil {
        t.Error(err)
        return
    }
    if len(records) != 1 {
        t.Error("query failed")
        return
    }
    err = m.DeleteRecord(record.Id)
    if err != nil {
        t.Error(err)
        return
    }
    newCount, err := m.Count()
    if err != nil {
        t.Error(err)
        return
    }
    if newCount-count != 0 {
        t.Error("delete failed")
        return
    }
}

func TestDeleteAll(t *testing.T) {
    TestAdd(t)
    var m RecordManager
    m = &RecordManagerSqlite3{}
    err := m.Open()
    if err != nil {
        t.Error(err)
        return
    }
    defer m.Close()
    err = m.DeleteAllRecord()
    if err != nil {
        t.Error(err)
        return
    }
    count, err := m.Count()
    if err != nil {
        t.Error(err)
        return
    }
    if count != 0 {
        t.Error("delete all failed")
        return
    }
}
