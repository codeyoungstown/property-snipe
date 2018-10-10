package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

const AUDITOR_URL = "http://oh-mahoning-auditor.publicaccessnow.com/DesktopModules/OWS/IM.aspx?&mpropertynumber=%s&_OWS_=lxC:1,lxP:0,s:1x,m:427,pm:116,p:89,lxSrc:dnn$ctr427$PropertyInfo,key:Module.Load,file:/DesktopModules/PropertyInfo/App_LocalResources/PropertyInfo.ascx.resx,pp:0"

func main() {
	arg := os.Args

	if len(arg) == 1 {
		sync()
		fmt.Println("Sync and compare done.")
		os.Exit(0)
	}

	cmd := arg[1]

	input := "" //FIXME
	if len(arg) > 2 {
		input = arg[2]
	}

	if cmd == "add" {
		add(input)
	} else if cmd == "remove" {
		remove(input)
	} else if cmd == "list" {
		list()
	} else {
		fmt.Fprintf(os.Stderr, "error: invalid command\n")
		os.Exit(1)
	}
}

func sync() {
	db, err := sql.Open("sqlite3", "./db.db")
	_checkErr(err)

	rows, _ := db.Query("SELECT id, parcel_id, owner FROM properties")

	var id int
	var parcel_id string
	var owner string

	updates := make(map[int]string)

	for rows.Next() {
		rows.Scan(&id, &parcel_id, &owner)

		new_owner := _get_owner(parcel_id)

		if owner != new_owner {
			fmt.Printf("New owner for %s: %s\n", parcel_id, new_owner)
			updates[id] = new_owner
		}
	}

	for k, v := range updates {
		statement, _ := db.Prepare("UPDATE properties SET owner=? WHERE id=?")
		statement.Exec(v, k)
	}

	db.Close()
}

func list() {
	db, err := sql.Open("sqlite3", "./db.db")
	_checkErr(err)

	rows, _ := db.Query("SELECT id, parcel_id, owner FROM properties")

	var id int
	var parcel_id string
	var owner string

	for rows.Next() {
		rows.Scan(&id, &parcel_id, &owner)
		fmt.Println(strconv.Itoa(id) + " " + parcel_id + " " + owner)
	}
	db.Close()
}

func add(parcel_id string) {
	db, err := sql.Open("sqlite3", "./db.db")
	_checkErr(err)

	// FIXME Check exists already and bail

	owner := _get_owner(parcel_id)

	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS properties (id INTEGER PRIMARY KEY, parcel_id TEXT, owner TEXT)")
	statement.Exec()

	statement, err = db.Prepare("INSERT INTO properties(parcel_id, owner) values(?,?)")
	_checkErr(err)

	statement.Exec(parcel_id, owner)

	db.Close()
}

func remove(parcel_id string) {
	db, err := sql.Open("sqlite3", "./db.db")
	_checkErr(err)

	statement, err := db.Prepare("DELETE FROM properties WHERE parcel_id=?")
	_checkErr(err)

	statement.Exec(parcel_id)

	db.Close()
}

func _checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func _get_owner(parcel_id string) string {
	url := fmt.Sprintf(AUDITOR_URL, parcel_id)
	res, _ := http.Get(url)
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	r, _ := regexp.Compile("Owner Name</td><td >([\\w ]+)</td>")

	// FIXME handle no match
	return r.FindStringSubmatch(string(body[:]))[1]
}
