package main

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type Row struct {
	Start int64  `json:"start"`
	Host  string `json:"host"`
	Fname string `json:"fname"`
}

func InitDb(db *sql.DB) (err error) {
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			host  TEXT    NOT NULL,
			start INTEGER NOT NULL,
			fname TEXT    NOT NULL
		)
	`)
	if err != nil {
		return
	}

	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx1 ON sessions(start)
	`)
	if err != nil {
		return
	}

	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx1 ON sessions(host)
	`)
	if err != nil {
		return
	}

	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx2 ON sessions(host, start)
	`)
	if err != nil {
		return
	}

	return
}

func InsertSession(db *sql.DB, host string, start int64, fname string) error {
	_, err := sq.Insert("sessions").
		Columns("host", "start", "fname").
		Values(host, start, fname).
		RunWith(db).Exec()
	return err
}

func GetSessions(db *sql.DB, page, limit int64) ([]Row, error) {
	rows, err := sq.Select("host", "start", "fname").
		From("sessions").
		OrderBy("start desc").
		Limit(uint64(limit)).
		Offset(uint64((page - 1) * limit)).
		RunWith(db).Query()
	if err != nil {
		return nil, err
	}

	data := make([]Row, 0, limit)

	for rows.Next() {
		d := Row{}
		err = rows.Scan(&d.Host, &d.Start, &d.Fname)
		if err != nil {
			return nil, err
		}
		data = append(data, d)
	}

	return data, nil
}
