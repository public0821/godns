package db

import (
	"fmt"
	"github.com/public0821/dnserver/errors"
)

type SysOption struct {
	Name  string
	Value string
}

func (m *SysOption) FuzzyQuery(dbm *DBManager, start, end int) (records []interface{}, err error) {
	err = errors.New("unimplemented")
	return
}
func (m *SysOption) Query(dbm *DBManager, start, offset int) (records []interface{}, err error) {
	if start < 0 || offset < 0 {
		err = errors.New("invalid arguments")
	}
	sql := "select name, value from sysoption "
	if len(m.Name) != 0 {
		sql += fmt.Sprintf(" where name='%s'", m.Name)
	}
	rows, err := dbm.db.Query(sql)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var option SysOption
		rows.Scan(&option.Name, &option.Value)
		records = append(records, option)
	}

	return
}

func (m *SysOption) Modify(dbm *DBManager) (err error) {
	sql := fmt.Sprintf("update sysoption set value='%s' where name='%s'",
		m.Value, m.Name)
	_, err = dbm.db.Exec(sql)
	return
}

func (m *SysOption) DeleteAll(dbm *DBManager) (err error) {
	err = errors.New("unimplemented")
	return
}

func (m *SysOption) Delete(dbm *DBManager) (err error) {
	err = errors.New("unimplemented")
	return
}

func (m *SysOption) Count(dbm *DBManager) (count int, err error) {
	err = errors.New("unimplemented")
	return
}

func (m *SysOption) Add(dbm *DBManager) (err error) {
	tx, err := dbm.db.Begin()
	if err != nil {
		return
	}
	stmt, err := tx.Prepare(`insert into sysoption(name,  value)
        values (?,?)`)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(m.Name, m.Value)
	if err != nil {
		return
	}
	tx.Commit()

	return
}
