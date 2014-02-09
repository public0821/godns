package db

import (
	"fmt"
	"github.com/public0821/dnserver/dns"
	"github.com/public0821/dnserver/errors"
	"net"
	"strings"
)

type RRecord struct {
	Id    int64
	Name  string
	Class uint16
	Type  uint16
	Value string
	Ttl   uint32
}

func (m *RRecord) isValidFormat() (err error) {
	if len(m.Name) == 0 || len(m.Value) == 0 {
		err = errors.New("name and value field should be specified")
		return
	}
	if m.Type == dns.TYPE_A || m.Type == dns.TYPE_MX || m.Type == dns.TYPE_AAAA {
		if net.ParseIP(m.Value) == nil {
			err = errors.New("invalid ip address")
			return
		}
	}
	if m.Type == dns.TYPE_PTR {
		if net.ParseIP(m.Name) == nil {
			err = errors.New("invalid ip address")
			return
		}
	}
	return
}

func (m *RRecord) DeleteAll(dbm *DBManager) (err error) {
	sql := `delete from record`
	_, err = dbm.db.Exec(sql)
	return
}

func (m *RRecord) Delete(dbm *DBManager) (err error) {
	sql := fmt.Sprintf(`delete from record where id=%d`, m.Id)
	_, err = dbm.db.Exec(sql)
	return
}

//TODO: refine code to avoid SQL injection
func (m *RRecord) Modify(dbm *DBManager) (err error) {
	err = m.isValidFormat()
	if err != nil {
		return
	}
	sql := fmt.Sprintf("update record set name='%s', class=%d, type=%d, value='%s', ttl=%d  where id=%d",
		m.Name, m.Class, m.Type, m.Value, m.Ttl, m.Id)
	_, err = dbm.db.Exec(sql)
	return
}

func getQuerySql(m *RRecord, fuzzy bool) string {
	sql := "select id, name, class, type, value, ttl from record "
	var conditions []string
	if m.Id != 0 {
		conditions = append(conditions, fmt.Sprintf(" id=%d ", m.Id))
	}
	if m.Class != 0 {
		conditions = append(conditions, fmt.Sprintf(" class=%d ", m.Class))
	}
	if m.Type != 0 {
		conditions = append(conditions, fmt.Sprintf(" type=%d ", m.Type))
	}
	if len(m.Name) != 0 {
		if fuzzy {
			conditions = append(conditions, fmt.Sprintf(" name like '%%%s%%' ", m.Name))
		} else {
			conditions = append(conditions, fmt.Sprintf(" name = '%s' ", m.Name))
		}
	}
	if len(m.Value) != 0 {
		if fuzzy {
			conditions = append(conditions, fmt.Sprintf(" value like '%%%s%%' ", m.Value))
		} else {
			conditions = append(conditions, fmt.Sprintf(" value = '%s' ", m.Value))
		}
	}
	if m.Ttl != 0 {
		conditions = append(conditions, fmt.Sprintf(" ttl=%d ", m.Ttl))
	}
	if len(conditions) > 0 {
		sql += " where " + strings.Join(conditions, "and")
	}
	sql += " order by name, class, type "

	return sql
}

//TODO: refine code to avoid SQL injection
func (m *RRecord) Query(dbm *DBManager, start, offset int) (records []interface{}, err error) {
	records, err = m.query(dbm, false, start, offset)
	return
}

func (m *RRecord) FuzzyQuery(dbm *DBManager, start, offset int) (records []interface{}, err error) {
	records, err = m.query(dbm, true, start, offset)
	return
}

func (m *RRecord) query(dbm *DBManager, fuzzy bool, start, offset int) (records []interface{}, err error) {
	if start < 0 || offset < 0 {
		err = errors.New("invalid arguments")
	}
	sql := getQuerySql(m, fuzzy)
	if offset != 0 {
		sql += fmt.Sprintf(" limit %d, %d ", start, offset)
	}
	fmt.Println(sql)
	rows, err := dbm.db.Query(sql)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var record RRecord
		rows.Scan(&record.Id, &record.Name, &record.Class, &record.Type, &record.Value, &record.Ttl)
		records = append(records, record)
	}

	return
}

func (m *RRecord) Count(dbm *DBManager) (count int, err error) {
	sql := "select count(*) from record "
	rows, err := dbm.db.Query(sql)
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

func (m *RRecord) Add(dbm *DBManager) (err error) {
	err = m.isValidFormat()
	if err != nil {
		return
	}
	tx, err := dbm.db.Begin()
	if err != nil {
		return
	}
	stmt, err := tx.Prepare(`insert into record(name, class, type, value, ttl)
        values (?,?,?,?,?)`)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(m.Name, m.Class, m.Type, m.Value, m.Ttl)
	if err != nil {
		return
	}
	tx.Commit()

	return
}
