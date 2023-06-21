package db

import "database/sql"

const (
	columnName = "name"
	columnBPM  = "bpm"
	columnTime = "time"

	tableData = "zoll_data"
)

type DBConn struct {
	Conn *sql.DB
}

func NewDBConnFromExisting(conn *sql.DB) *DBConn {
	return &DBConn{
		Conn: conn,
	}
}
