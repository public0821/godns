package db

import (
    "database/sql"
    "github.com/coopernurse/gorp"
    _ "github.com/mattn/go-sqlite3"
    "github.com/public0821/dnserver/util"
    "log"
    "os"
)

const db_path = "/etc/dnserver"
const db_file = db_path + "/db.sqlite3"

func InitDb() (err error) {
    //file not exist
    if _, err = os.Stat(db_file); os.IsNotExist(err) {
        err = os.MkdirAll(db_path, 0744)
        if err != nil {
            return
        }
        err = initDb()
        if err != nil {
            os.Remove(db_file)
            return
        }
    }
    return
}

func initDb() (err error) {
    // connect to db using standard Go database/sql API
    db, err := sql.Open("sqlite3", db_file)
    if err != nil {
        return
    }
    defer db.Close()

    // construct a gorp DbMap
    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

    // add a table, setting the table name to 'posts' and
    // specifying that the Id property is an auto incrementing PK
    dbmap.AddTableWithName(RRecord{}, "rrecord").SetKeys(true, "Id")
    dbmap.AddTableWithName(SysOption{}, "sysoption").SetKeys(true, "Id")
    dbmap.AddTableWithName(ForwardServer{}, "forward_server").SetKeys(true, "Id")
    dbmap.AddTableWithName(User{}, "user").SetKeys(true, "Id")

    // create the table. in a production system you'd generally
    // use a migration tool, or create the tables via scripts
    err = dbmap.CreateTables()
    if err != nil {
        return
    }

    //init data
    var option SysOption
    option.Name = "mode"
    option.Value = "forward"
    err = dbmap.Insert(&option)
    if err != nil {
        return
    }
    var servers []string
    servers, err = util.GetDnsServer()
    if err != nil {
        return
    }
    log.Println("dnsserver:", servers)
    for _, server := range servers {
        fs := ForwardServer{Ip: server}
        err = dbmap.Insert(&fs)
        if err != nil {
            return
        }
    }
    var user User
    user.Name = "admin"
    user.Pwd = "admin"
    err = dbmap.Insert(&user)
    if err != nil {
        return
    }

    return
}

func GetDbmap(dbmap *gorp.DbMap, err error) {
    db, err := sql.Open("sqlite3", db_file)
    if err != nil {
        return
    }

    // construct a gorp DbMap
    dbmap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
    return
}
