package db

import (
    "os"
    "testing"
)

func TestInitDb(t *testing.T) {
    os.Remove(DB_FILE)
    err := InitDb()
    if err != nil {
        t.Error(err)
        return
    }
    err = InitDb()
    if err != nil {
        t.Error(err)
        return
    }
}
