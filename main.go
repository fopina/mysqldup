package main

/*
https://gist.github.com/csonuryilmaz/3f8f92fdad007f97986e61ad79aeb514
*/

import (
	"database/sql"
	"fmt"
	"os"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	flag "github.com/spf13/pflag"
	"golang.org/x/crypto/ssh/terminal"
)

var version string = "DEV"
var date string

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	hostnamePtr := flag.StringP("hostname", "h", "127.0.0.1", "connect to host")
	portPtr := flag.IntP("port", "P", 3306, "port number to use for connection")
	userPtr := flag.StringP("user", "u", "root", "user for login")
	usePasswordPtr := flag.BoolP("password", "p", false, "use password when connecting to server (read from tty)")
	forcePtr := flag.BoolP("force", "f", false, "drop NEW_DB_NAME if it already exists")
	versionPtr := flag.BoolP("version", "V", false, "output version information and exit")
	helpPtr := flag.Bool("help", false, "this screen")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] OLD_DB_NAME NEW_DB_NAME\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *versionPtr {
		fmt.Println("Version: " + version + " (built on " + date + ")")
		return
	}

	if *helpPtr {
		flag.Usage()
		return
	}

	if flag.NArg() < 2 {
		flag.Usage()
		fmt.Fprintln(os.Stderr, "missing mandatory positional arguments")
		os.Exit(2)
	}

	oldDB := flag.Arg(0)
	newDB := flag.Arg(1)
	password := ""
	if *usePasswordPtr {
		fmt.Print("Enter Password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		check(err)
		password = string(bytePassword)
	}

	con, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/", *userPtr, password, *hostnamePtr, *portPtr))
	check(err)
	defer con.Close()

	_, err = con.Exec("USE " + oldDB)
	check(err)
	var dbCharset, dbCollation string
	r := con.QueryRow("SELECT @@character_set_database, @@collation_database;")
	r.Scan(&dbCharset, &dbCollation)

	if *forcePtr {
		_, err = con.Exec("DROP DATABASE IF EXISTS " + newDB)
		check(err)
	}

	_, err = con.Exec("CREATE DATABASE " + newDB + " CHARACTER SET " + dbCharset + " COLLATE " + dbCollation)
	check(err)

	rows, err := con.Query("SHOW FULL TABLES WHERE Table_Type = 'BASE TABLE'")
	check(err)
	defer rows.Close()
	var tablesToClone []string
	for rows.Next() {
		var (
			tableName, tableType string
		)
		err = rows.Scan(&tableName, &tableType)
		check(err)
		tablesToClone = append(tablesToClone, tableName)
	}
	fmt.Printf("%v tables to clone\n", len(tablesToClone))

	_, err = con.Exec("USE " + newDB)
	check(err)

	_, err = con.Exec("set foreign_key_checks = 0")
	check(err)

	for tableInd, table := range tablesToClone {
		fmt.Printf("[%d / %d] cloning %v\n", tableInd+1, len(tablesToClone), table)
		_, err = con.Exec("CREATE TABLE " + table + " SELECT * FROM " + oldDB + "." + table)
		check(err)
	}

	_, err = con.Exec("set foreign_key_checks = 1")
	check(err)
}
