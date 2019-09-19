package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"
)

// SQLDatabase struct to hold the connection
type SQLDatabase struct {
	Connection *sql.DB
}

// Results contains the column names and the data
type Results struct {
	Columns []string
	Data    [][]string
}

var sqlDb SQLDatabase

// GetConnection Make a connection to the database to execute queries against
func GetConnection() (bool, error) {
	if sqlDb.Connection != nil {
		err := sqlDb.Connection.Ping()
		if err != nil {
			return false, err
		}
		return true, nil
	}

	userid := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	server := os.Getenv("DATABASE_SERVER")
	database := os.Getenv("DATABASE_NAME")

	dsn := "server=" + server + ";user id=" + userid + ";password=" + password

	if len(database) != 0 {
		dsn += ";database=" + database
	}

	fmt.Println(dsn)

	connection, err := sql.Open("mssql", dsn)
	if err != nil {
		return false, err
	}

	err = connection.Ping()
	if err != nil {
		return false, err
	}

	sqlDb = SQLDatabase{Connection: connection}
	return true, nil
}

// ExecuteQuery Run the provided command and return the results in a 2 dimensional slice
// where slice[0] are column names, and all other slices are the resulting rows
func ExecuteQuery(cmd string) (*Results, error) {
	rows, err := sqlDb.Connection.Query(cmd)
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

	var results Results
	results.Columns = cols
	results.Data = make([][]string, 0)

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

		rowSlice := make([]string, len(cols))

		for i := 0; i < len(vals); i++ {
			rowSlice[i] = printValue(vals[i].(*interface{}))
		}
		rowCount++
		// Append the row data to the result slice (possibly not the most efficient way to do this...)
		results.Data = append(results.Data, rowSlice)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return &results, nil
}

// ExecuteWithTransaction opens a transaction and executes the cmd. Rollback happens if there is an error
func ExecuteWithTransaction(cmd string) error {
	tx, err := sqlDb.Connection.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	if _, err = tx.Exec(cmd); err != nil {
		return err
	}

	return nil
}

// ExecuteStatement executes the supplied procedure with parameters and returns results only if returnResults param is true
func ExecuteStatement(name string, executeOnly bool, parameters map[string]*interface{}) (*Results, error) {
	statement := "EXEC " + name + " "

	for key, value := range parameters {
		statement += "@" + key + " = " + printValue(value) + ", "
	}

	statement = strings.TrimRight(statement, ", ")

	if executeOnly {
		if _, err := sqlDb.Connection.Exec(statement); err != nil {
			return nil, err
		}
		return nil, nil
	}

	results, err := ExecuteQuery(statement)

	if err != nil {
		return nil, err
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
		return fmt.Sprint(v)
	}
}
