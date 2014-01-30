package db

import (
    //"fmt"
    "testing"
)

func TestUserAdd(t *testing.T) {
    count, err := Count(&User{})
    if err != nil {
        t.Error(err)
        return
    }

    var user User
    user.Name = "testadmin"
    user.Pwd = "password"
    err = Add(&user)
    if err != nil {
        t.Error(err)
        return
    }
    newCount, err := Count(&user)
    if err != nil {
        t.Error(err)
        return
    }
    if newCount-count != 1 {
        t.Error("add failed")
        return
    }
}

func TestUserQuery(t *testing.T) {
    var user User
    user.Name = "testadmin"
    records, err := Query(&user, 0, 0)
    if err != nil {
        t.Error(err)
        return
    }
    if len(records) != 1 {
        t.Error("query failed")
        return
    }
}

func TestUserModify(t *testing.T) {
    var user User
    user.Name = "testadmin"
    records, err := Query(&user, 0, 0)
    if err != nil {
        t.Error(err)
        return
    }
    if len(records) != 1 {
        t.Error("query failed")
        return
    }
    newUser, _ := records[0].(User)
    newUser.Name = "testadmin"
    newUser.Pwd = "password1"
    err = Modify(&newUser)
    if err != nil {
        t.Error(err)
        return
    }
    newRecords, err := Query(&newUser, 0, 0)
    if err != nil {
        t.Error(err)
        return
    }
    if len(newRecords) != 1 {
        t.Error("modify rrecord failed")
        return
    }

}
