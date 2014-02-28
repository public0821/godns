package db

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "github.com/public0821/dnserver/util"
    "log"
    "os"
    "sync"
)

const DB_PATH = "/etc/dnserver"
const DB_FILE = DB_PATH + "/db.sqlite3"

type DBManager struct {
    sync.RWMutex
    db  *sql.DB
}

var gdbm DBManager

func NewDBManager() (dbm *DBManager, err error) {
    if gdbm.db != nil {
        dbm = &gdbm
        return
    }

    gdbm.Lock()
    defer gdbm.Unlock()
    if gdbm.db != nil {
        dbm = &gdbm
        return
    }
    var tempdb *sql.DB
    //file not exist
    if _, err = os.Stat(DB_FILE); os.IsNotExist(err) {
        err = os.MkdirAll(DB_PATH, 0744)
        if err != nil {
            return
        }
        tempdb, err = sql.Open("sqlite3", DB_FILE)
        if err != nil {
            return
        }
        var tempDbm DBManager
        tempDbm.db = tempdb
        sql := `create table record (
        id integer not null primary key
        , name text
        , class integer
        , type integer
        , ttl integer
        , value text);`
        _, err = tempdb.Exec(sql)
        if err != nil {
            tempdb.Close()
            os.Remove(DB_FILE)
            return
        }
        sql = `create table sysoption (
        id integer not null primary key
        , name text
        , value text);`
        _, err = tempdb.Exec(sql)
        if err != nil {
            tempdb.Close()
            os.Remove(DB_FILE)
            return
        }
        var option SysOption
        option.Name = "mode"
        option.Value = "forward"
        err = option.Add(&tempDbm)
        if err != nil {
            tempdb.Close()
            os.Remove(DB_FILE)
            return
        }

        sql = `create table forward_server (
        id integer not null primary key
        , ip text);`
        _, err = tempdb.Exec(sql)
        if err != nil {
            tempdb.Close()
            os.Remove(DB_FILE)
            return
        }
        var servers []string
        servers, err = util.GetDnsServer()
        log.Println("dnsserver:", servers)
        for _, server := range servers {
            fs := ForwardServer{Ip: server}
            err = fs.Add(&tempDbm)
            if err != nil {
                tempdb.Close()
                os.Remove(DB_FILE)
                return
            }
        }

        sql = `create table user (
        id integer not null primary key
        , name text
        , pwd text);`
        _, err = tempdb.Exec(sql)
        if err != nil {
            tempdb.Close()
            os.Remove(DB_FILE)
            return
        }
        var user User
        user.Name = "admin"
        user.Pwd = "admin"
        err = user.Add(&tempDbm)
        if err != nil {
            tempdb.Close()
            os.Remove(DB_FILE)
            return
        }
    } else {
        tempdb, err = sql.Open("sqlite3", DB_FILE)
        if err != nil {
            return
        }
    }

    gdbm.db = tempdb
    dbm = &gdbm

    return
}

func DeleteAll(m Model) (err error) {
    dbm, err := NewDBManager()
    if err != nil {
        return
    }
    dbm.Lock()
    err = m.DeleteAll(dbm)
    dbm.Unlock()
    return
}
func Add(m Model) (err error) {
    dbm, err := NewDBManager()
    if err != nil {
        return
    }
    dbm.Lock()
    err = m.Add(dbm)
    dbm.Unlock()
    return
}
func Delete(m Model) (err error) {
    dbm, err := NewDBManager()
    if err != nil {
        return
    }
    dbm.Lock()
    err = m.Delete(dbm)
    dbm.Unlock()
    return
}
func Modify(m Model) (err error) {
    dbm, err := NewDBManager()
    if err != nil {
        return
    }
    dbm.Lock()
    err = m.Modify(dbm)
    dbm.Unlock()
    return
}

func Query(m Model, start, offset int) (records []interface{}, err error) {
    dbm, err := NewDBManager()
    if err != nil {
        return
    }
    dbm.Lock()
    records, err = m.Query(dbm, start, offset)
    dbm.Unlock()
    return
}
func FuzzyQuery(m Model, start, offset int) (records []interface{}, err error) {
    dbm, err := NewDBManager()
    if err != nil {
        return
    }
    dbm.Lock()
    records, err = m.FuzzyQuery(dbm, start, offset)
    dbm.Unlock()
    return
}
func Count(m Model) (count int, err error) {
    dbm, err := NewDBManager()
    if err != nil {
        return
    }
    dbm.Lock()
    count, err = m.Count(dbm)
    dbm.Unlock()
    return
}
