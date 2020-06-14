package main

import (
	"log"

	"gopkg.in/mgo.v2"
)

var db *mgo.Session

func dialdb() error {
	var err error
	log.Println("dialing mongodb: localhost")
	db, err = mgo.Dial("localhost") // TODO: make this configurable
	return err
}

func closeddb() {
	db.Close()
	log.Println("closed database connection")
}

func main() {

}
