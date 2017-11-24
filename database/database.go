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
	userid := "sa"
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

	return connection, err
}

func Execute(db *sql.DB, cmd string) ([][]string, error) {
	rows, err := db.Query(cmd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	if cols == nil {
		return nil, nil
	}

	// TODO: Somehow need to dynamically allocate the slice size based on the amount of rows being returned
	results := make([][]string, 3)
	results[0] = cols

	vals := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		vals[i] = new(interface{})
	}

	rowCount := 1

	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			fmt.Println(err)
			continue
		}

		results[rowCount] = append(make([]string, len(cols)))

		for i := 0; i < len(vals); i++ {
			results[rowCount][i] = printValue(vals[i].(*interface{}))
		}
		rowCount++
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return results, nil
}

func printValue(pval *interface{}) string {
	switch v := (*pval).(type) {
	case nil:
		return "NULL"
	case bool:
		if v {
			return "1"
		}
		return "0"
	case []byte:
		return string(v)
	case time.Time:
		return v.Format("2006-01-02 15:04:05.999")
	default:
		if str, ok := v.(string); ok {
			return string(str)
		}
		return ""
	}
}
