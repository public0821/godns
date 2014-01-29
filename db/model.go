package db

import (
//"github.com/public0821/dnserver/errors"
)

type Model interface {
    DeleteAll(dbm *DBManager) (err error)
    Add(dbm *DBManager) (err error)

    //delete the record whoes id is equal to this record's id
    Delete(dbm *DBManager) (err error)

    Modify(dbm *DBManager) (err error)
    Query(dbm *DBManager, start, end int) (records []interface{}, err error)

    //return the count of all record in table
    Count(dbm *DBManager) (count int, err error)
}
