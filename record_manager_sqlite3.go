package main

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "os"
)

type Record struct {
    Id    int64
    Name  string
    Class uint8
    Type  uint8
    Value string
}

const DB_FILE = "./record.db"

func initDB() (err error) {
    db, err := sql.Open("sqlite3", DB_FILE)
    if err != nil {
        return
    }
    defer db.Close()

    sql := `create table record (
        id integer not null primary key
        , name text
        , class integer
        , type integer
        , value text);`
    _, err = db.Exec(sql)
    if err != nil {
        return
    }
    return
}

func deleteAllRecord() (err error) {
    db, err := getDB()
    if err != nil {
        return
    }
    sql := `delete from record`
    _, err = db.Exec(sql)
    if err != nil {
        return
    }
    return
}

var gdb *sql.DB = nil

func getDB() (db *sql.DB, err error) {
    if gdb != nil {
        db = gdb
        return
    }
    //file not exist
    if _, err = os.Stat(DB_FILE); os.IsNotExist(err) {
        err = initDB()
        if err != nil {
            return
        }
    }

    db, err = sql.Open("sqlite3", DB_FILE)
    if err != nil {
        return
    }
    gdb = db
    return
}

func getRecord() (records []Record, err error) {
    db, err := getDB()
    if err != nil {
        return
    }

    rows, err := db.Query("select id, name, class, type, value from record order by name")
    if err != nil {
        return
    }
    defer rows.Close()
    for rows.Next() {
        var record Record
        rows.Scan(&record.Id, &record.Name, &record.Class, &record.Type, &record.Value)
        records = append(records, record)
    }
    rows.Close()

    return
}

func addRecord(record *Record) (err error) {
    db, err := getDB()
    if err != nil {
        return
    }
    tx, err := db.Begin()
    if err != nil {
        return
    }
    stmt, err := tx.Prepare(`insert into record(name, class, type, value)
        values (?,?,?,?)`)
    if err != nil {
        return
    }
    defer stmt.Close()
    _, err = stmt.Exec(record.Name, record.Class, record.Type, record.Value)
    if err != nil {
        return
    }
    tx.Commit()

    return
}

func main() {
    var record Record
    record.Name = "www.test.com"
    record.Class = 1
    record.Type = 1
    record.Value = "10.32.171.60"
    err := deleteAllRecord()
    if err != nil {
        fmt.Println(err)
        return
    }
    for i := 0; i < 10000; i++ {
        err = addRecord(&record)
        if err != nil {
            fmt.Println(err)
            return
        }
    }
    records, err := getRecord()
    if err != nil {
        fmt.Println(err)
        return
    } else {
        fmt.Println(len(records))
    }

}
