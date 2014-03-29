package db

import (
    "database/sql"
    "github.com/coopernurse/gorp"
    _ "github.com/mattn/go-sqlite3"
    "github.com/public0821/dnserver/util"
    "log"
    "os"
)

const DB_PATH = "/etc/dnserver"
const DB_FILE = DB_PATH + "/db.sqlite3"

func InitDb() (err error) {
    //file not exist
    if _, err = os.Stat(DB_FILE); os.IsNotExist(err) {
        err = os.MkdirAll(DB_PATH, 0744)
        if err != nil {
            return
        }
        err = initDb()
        if err != nil {
            os.Remove(DB_FILE)
            return
        }
    }

    return
}

func addTable(dbmap *gorp.DbMap) {
    // add a table, setting the table name to 'posts' and
    // specifying that the Id property is an auto incrementing PK
    dbmap.AddTableWithName(Record{}, "record").SetKeys(true, "Id")
    dbmap.AddTableWithName(SysOption{}, "sysoption").SetKeys(false, "Name")
    dbmap.AddTableWithName(ForwardServer{}, "forward_server").SetKeys(false, "Ip")
    dbmap.AddTableWithName(User{}, "user").SetKeys(false, "Name")
}

func initDb() (err error) {
    // connect to db using standard Go database/sql API
    db, err := sql.Open("sqlite3", DB_FILE)
    if err != nil {
        return
    }
    defer db.Close()

    // construct a gorp DbMap
    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

    addTable(dbmap)

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

func OpenDbmap() (dbmap *gorp.DbMap, err error) {
    db, err := sql.Open("sqlite3", DB_FILE)
    if err != nil {
        return
    }

    // construct a gorp DbMap
    dbmap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
    addTable(dbmap)
    return
}
func CloseDbmap(dbmap *gorp.DbMap) {
    dbmap.Db.Close()
}
