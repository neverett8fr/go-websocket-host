package db

import (
	"fmt"
	"log"
	"time"
)

func (conn *DBConn) CreateEntry(data string) error {

	_, err := conn.Conn.Exec(
		fmt.Sprintf(
			"INSERT INTO %s(%s, %s, %s) VALUES($1, $2, $3)",
			tableData, columnName, columnBPM, columnTime,
		),
		data, 80, time.Now(),
	)
	if err != nil {
		err := fmt.Errorf("error inserting new entry, err %v", err)
		log.Println(err)
		return err
	}

	return nil
}
