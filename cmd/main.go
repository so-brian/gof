package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"os"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	// dsn := "server=" + *server + ";user id=" + *userid + ";password=" + *password + ";database=" + *database
	dsn := "Server=tcp:sobrian.database.windows.net,1433;Initial Catalog=SoData;Persist Security Info=False;User ID=username;Password={password};MultipleActiveResultSets=False;Encrypt=True;TrustServerCertificate=False;Connection Timeout=30;"
	db, err := sql.Open("mssql", dsn)
	if err != nil {
		fmt.Println("Cannot connect: ", err.Error())
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("Cannot connect: ", err.Error())
		return
	}
	r := bufio.NewReader(os.Stdin)
	for {
		_, err = os.Stdout.Write([]byte("> "))
		if err != nil {
			fmt.Println(err)
			return
		}
		cmd, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				return
			}
			fmt.Println(err)
			return
		}
		err = exec(db, cmd)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func exec(db *sql.DB, cmd string) error {
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
