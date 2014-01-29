package db

import (
    //"fmt"
    "testing"
)

func TestAdd(t *testing.T) {
    count, err := Count(&RRecord{})
    if err != nil {
        t.Error(err)
        return
    }

    var rrecord RRecord
    rrecord.Name = "www.test.com"
    rrecord.Class = 1
    rrecord.Type = 1
    rrecord.Value = "10.32.171.60"
    err = Add(&rrecord)
    if err != nil {
        t.Error(err)
        return
    }
    newCount, err := Count(&rrecord)
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
    var rrecord RRecord
    rrecord.Name = "www.test.com"
    rrecord.Class = 1
    rrecord.Type = 1
    rrecord.Value = "10.32.171.60"
    records, err := Query(&rrecord, 0, 0)
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
    var rrecord RRecord
    rrecord.Name = "www.test.com"
    rrecord.Class = 1
    rrecord.Type = 1
    rrecord.Value = "10.32.171.60"
    records, err := Query(&rrecord, 0, 0)
    if err != nil {
        t.Error(err)
        return
    }
    if len(records) != 1 {
        t.Error("query failed")
        return
    }
    newRecord, _ := records[0].(RRecord)
    newRecord.Name = "www.test1.com"
    newRecord.Class = 2
    newRecord.Type = 2
    newRecord.Value = "11.32.171.60"
    err = Modify(&newRecord)
    if err != nil {
        t.Error(err)
        return
    }
    newRecords, err := Query(&newRecord, 0, 0)
    if err != nil {
        t.Error(err)
        return
    }
    if len(newRecords) != 1 {
        t.Error("modify rrecord failed")
        return
    }

}

func TestDelete(t *testing.T) {
    TestAdd(t)
    count, err := Count(&RRecord{})
    if err != nil {
        t.Error(err)
        return
    }
    var rrecord RRecord
    rrecord.Name = "www.test.com"
    rrecord.Class = 1
    rrecord.Type = 1
    rrecord.Value = "10.32.171.60"
    records, err := Query(&rrecord, 0, 0)
    if err != nil {
        t.Error(err)
        return
    }
    if len(records) != 1 {
        t.Error("query failed")
        return
    }
    err = Delete(&rrecord)
    if err != nil {
        t.Error(err)
        return
    }
    newCount, err := Count(&RRecord{})
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
    err := DeleteAll(&RRecord{})
    if err != nil {
        t.Error(err)
        return
    }
    count, err := Count(&RRecord{})
    if err != nil {
        t.Error(err)
        return
    }
    if count != 0 {
        t.Error("delete all failed")
        return
    }
}
