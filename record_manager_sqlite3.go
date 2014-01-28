package main

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "os"
    "strings"
)

const DB_FILE = "./record.db"

type RecordManagerSqlite3 struct {
    db *sql.DB
}

func (m *RecordManagerSqlite3) Open() (err error) {
    //file not exist
    if _, err = os.Stat(DB_FILE); os.IsNotExist(err) {
        m.db, err = sql.Open("sqlite3", DB_FILE)
        if err != nil {
            return
        }
        sql := `create table record (
        id integer not null primary key
        , name text
        , class integer
        , type integer
        , ttl integer
        , value text);`
        _, err = m.db.Exec(sql)
        if err != nil {
            return
        }
    } else {
        m.db, err = sql.Open("sqlite3", DB_FILE)
        if err != nil {
            return
        }
    }
    return
}

func (m *RecordManagerSqlite3) Close() (err error) {
    err = m.db.Close()
    return
}

func (m *RecordManagerSqlite3) DeleteAllRecord() (err error) {
    sql := `delete from record`
    _, err = m.db.Exec(sql)
    return
}

func (m *RecordManagerSqlite3) DeleteRecord(id int64) (err error) {
    sql := fmt.Sprintf(`delete from record where id=%d`, id)
    _, err = m.db.Exec(sql)
    return
}

//TODO: refine code to avoid SQL injection
func (m *RecordManagerSqlite3) ModifyRecord(record *Record) (err error) {
    sql := fmt.Sprintf("update record set name='%s', class=%d, type=%d, value='%s', ttl=%d  where id=%d",
        record.Name, record.Class, record.Type, record.Value, record.Ttl, record.Id)
    _, err = m.db.Exec(sql)
    return
}

//TODO: refine code to avoid SQL injection
func (m *RecordManagerSqlite3) QueryRecord(record *Record) (records []Record, err error) {
    sql := "select id, name, class, type, value, ttl from record "
    var conditions []string
    if record.Id != 0 {
        conditions = append(conditions, fmt.Sprintf(" id=%d ", record.Id))
    }
    if record.Class != 0 {
        conditions = append(conditions, fmt.Sprintf(" class=%d ", record.Class))
    }
    if record.Type != 0 {
        conditions = append(conditions, fmt.Sprintf(" type=%d ", record.Type))
    }
    if len(record.Name) != 0 {
        conditions = append(conditions, fmt.Sprintf(" name='%s' ", record.Name))
    }
    if len(record.Value) != 0 {
        conditions = append(conditions, fmt.Sprintf(" value='%s' ", record.Value))
    }
    if record.Ttl != 0 {
        conditions = append(conditions, fmt.Sprintf(" ttl=%d ", record.Ttl))
    }
    if len(conditions) > 0 {
        sql += " where " + strings.Join(conditions, "and")
    }
    sql += " order by name, class, type "
    rows, err := m.db.Query(sql)
    if err != nil {
        return
    }
    defer rows.Close()
    for rows.Next() {
        var record Record
        rows.Scan(&record.Id, &record.Name, &record.Class, &record.Type, &record.Value, &record.Ttl)
        records = append(records, record)
    }

    return
}

func (m *RecordManagerSqlite3) Count() (count int, err error) {
    sql := "select count(*) from record "
    rows, err := m.db.Query(sql)
    if err != nil {
        return
    }
    defer rows.Close()
    for rows.Next() {
        rows.Scan(&count)
        return
    }

    return
}

func (m *RecordManagerSqlite3) AddRecord(record *Record) (err error) {
    tx, err := m.db.Begin()
    if err != nil {
        return
    }
    stmt, err := tx.Prepare(`insert into record(name, class, type, value, ttl)
        values (?,?,?,?,?)`)
    if err != nil {
        return
    }
    defer stmt.Close()
    _, err = stmt.Exec(record.Name, record.Class, record.Type, record.Value, record.Ttl)
    if err != nil {
        return
    }
    tx.Commit()

    return
}
