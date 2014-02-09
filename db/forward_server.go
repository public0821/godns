package db

import (
	"fmt"
	"github.com/public0821/dnserver/errors"
	"net"
)

type ForwardServer struct {
	Ip string
}

func (m *ForwardServer) FuzzyQuery(dbm *DBManager, start, end int) (records []interface{}, err error) {
	err = errors.New("unimplemented")
	return
}

func (m *ForwardServer) Modify(dbm *DBManager) (err error) {
	err = errors.New("unimplemented")
	return
}

func (m *ForwardServer) Count(dbm *DBManager) (count int, err error) {
	err = errors.New("unimplemented")
	return
}

func (m *ForwardServer) Add(dbm *DBManager) (err error) {
	if net.ParseIP(m.Ip) == nil {
		err = errors.New("invalid ip address")
		return
	}
	tx, err := dbm.db.Begin()
	if err != nil {
		return
	}
	stmt, err := tx.Prepare(`insert into forward_server(ip)
        values (?)`)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(m.Ip)
	if err != nil {
		return
	}
	tx.Commit()

	return
}

func (m *ForwardServer) Query(dbm *DBManager, start, offset int) (records []interface{}, err error) {
	if start < 0 || offset < 0 {
		err = errors.New("invalid arguments")
	}
	sql := "select ip from forward_server "
	rows, err := dbm.db.Query(sql)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var server ForwardServer
		rows.Scan(&server.Ip)
		records = append(records, server)
	}

	return
}

func (m *ForwardServer) DeleteAll(dbm *DBManager) (err error) {
	sql := `delete from forward_server`
	_, err = dbm.db.Exec(sql)
	return
}

func (m *ForwardServer) Delete(dbm *DBManager) (err error) {
	sql := fmt.Sprintf(`delete from forward_server where ip='%s'`, m.Ip)
	_, err = dbm.db.Exec(sql)
	return
}
