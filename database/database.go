package database

import (
	"database/sql"
	"fmt"
	"time"
)

type SqlDatabase struct {
	Connection *sql.DB
}

func GetConnection() (*sql.DB, error) {
	userid := "username"
	password := "p@ssw0rd"
	server := "localhost"

	dsn := "server=" + server + ";user id=" + userid + ";password=" + password // + ";database=" + database

	connection, err := sql.Open("mssql", dsn)
	if err != nil {
		fmt.Println("Cannot connect: ", err.Error())
		return nil, err
	}
	err = connection.Ping()
	if err != nil {
		fmt.Println("Cannot connect: ", err.Error())
		return nil, err
	}
	// defer db.Close()

	return connection, err
}

func Execute(db *sql.DB, cmd string) error {
	rows, err := db.Query(cmd)
	if err != nil {
		return err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	if cols == nil {
		return nil
	}
	vals := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		vals[i] = new(interface{})
		if i != 0 {
			fmt.Print("\t")
		}
		fmt.Print(cols[i])
	}
	fmt.Println()
	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for i := 0; i < len(vals); i++ {
			if i != 0 {
				fmt.Print("\t")
			}
			printValue(vals[i].(*interface{}))
		}
		fmt.Println()

	}
	if rows.Err() != nil {
		return rows.Err()
	}
	return nil
}

func printValue(pval *interface{}) {
	switch v := (*pval).(type) {
	case nil:
		fmt.Print("NULL")
	case bool:
		if v {
			fmt.Print("1")
		} else {
			fmt.Print("0")
		}
	case []byte:
		fmt.Print(string(v))
	case time.Time:
		fmt.Print(v.Format("2006-01-02 15:04:05.999"))
	default:
		fmt.Print(v)
	}
}
